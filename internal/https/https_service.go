package https

import (
	"bufio"
	"crypto/tls"
	"errors"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"
)

type HttpsService struct {
	proxyWriter             http.ResponseWriter
	interceptedHttpsRequest *http.Request
	proxyHttpsRequest       *http.Request
	scheme                  string
	config                  *tls.Config
}

func NewHttpsService(writer http.ResponseWriter, interceptedHttpsRequest *http.Request) *HttpsService {
	requestedUrl, _ := url.Parse(interceptedHttpsRequest.RequestURI)
	return &HttpsService{
		proxyWriter:             writer,
		interceptedHttpsRequest: interceptedHttpsRequest,
		scheme:                  requestedUrl.Scheme,
	}
}

func (hs *HttpsService) ProxyHttpsRequest() error {
	hijackedConn, err := hs.interceptConnection()
	if err != nil {
		return err
	}

	TCPClientConn, err := hs.initializeTCPClient(hijackedConn)
	if err != nil {
		return err
	}

	TCPServerConn, err := hs.initializeTCPServer()
	if err != nil {
		return err
	}

	err = hs.doHttpsRequest(TCPClientConn, TCPServerConn)
	if err != nil {
		return err
	}

	defer hijackedConn.Close()
	defer TCPClientConn.Close()

	return nil
}

func (hs *HttpsService) generateCertificate() (tls.Certificate, error) {
	rootDir, err := os.Getwd()
	if err != nil {
		return tls.Certificate{}, err
	}

	cmdGenDir := rootDir + "/genCerts"
	certsDir := cmdGenDir + "/certs/"

	certFilename := certsDir + hs.scheme + ".crt"

	_, errStat := os.Stat(certFilename)
	if os.IsNotExist(errStat) {
		genCommand := exec.Command(cmdGenDir+"/gen_cert.sh", hs.scheme, strconv.Itoa(rand.Intn(100000000)))

		_, err := genCommand.CombinedOutput()
		if err != nil {
			log.Fatal(err)
			return tls.Certificate{}, err
		}
	}

	cert, err := tls.LoadX509KeyPair(certFilename, cmdGenDir+"/cert.key")
	if err != nil {
		return tls.Certificate{}, err
	}

	return cert, nil
}

func (hs *HttpsService) interceptConnection() (net.Conn, error) {
	hijacker, ok := hs.proxyWriter.(http.Hijacker)
	if !ok {
		return nil, errors.New("creating hijacker failed")
	}

	hijackedConn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}

	_, err = hijackedConn.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))

	if err != nil {
		hijackedConn.Close()
		return nil, err
	}
	return hijackedConn, nil
}

func (hs *HttpsService) initializeTCPClient(hijackedConn net.Conn) (*tls.Conn, error) {
	cert, err := hs.generateCertificate()
	if err != nil {
		return nil, err
	}

	hs.config = &tls.Config{Certificates: []tls.Certificate{cert}, ServerName: hs.scheme}

	TCPClientConn := tls.Server(hijackedConn, hs.config)
	err = TCPClientConn.Handshake()
	if err != nil {
		TCPClientConn.Close()
		hijackedConn.Close()
		return nil, err
	}

	clientReader := bufio.NewReader(TCPClientConn)
	TCPClientRequest, err := http.ReadRequest(clientReader)
	if err != nil {
		return nil, err
	}

	hs.proxyHttpsRequest = TCPClientRequest

	return TCPClientConn, nil
}

func (hs *HttpsService) initializeTCPServer() (*tls.Conn, error) {
	TCPServerConn, err := tls.Dial("tcp", hs.interceptedHttpsRequest.Host, hs.config)
	if err != nil {
		return nil, err
	}

	return TCPServerConn, nil
}

func (hs *HttpsService) doHttpsRequest(TCPClientConn *tls.Conn, TCPServerConn *tls.Conn) error {
	rawReq, err := httputil.DumpRequest(hs.proxyHttpsRequest, true)
	_, err = TCPServerConn.Write(rawReq)
	if err != nil {
		return err
	}

	serverReader := bufio.NewReader(TCPServerConn)
	TCPServerResponse, err := http.ReadResponse(serverReader, hs.proxyHttpsRequest)
	if err != nil {
		return err
	}

	rawResp, err := httputil.DumpResponse(TCPServerResponse, true)
	_, err = TCPClientConn.Write(rawResp)
	if err != nil {
		return err
	}

	return nil
}
