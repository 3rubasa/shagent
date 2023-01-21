package sonoffr3rf

import (
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/3rubasa/shagent/controllers/relay/sonoffr3rf/mockosservicesprovider"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_GetState_On(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	l, err := net.Listen("tcp", "127.0.0.1:8081")
	assert.NoError(t, err)

	expectedState := "on"
	expectedURLPath := "/zeroconf/info"
	actualURLPath := ""
	expectedBody := `{ 
		"deviceid": "", 
		"data": { } 
	 }`
	actualBody := ""

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		actualURLPath = r.URL.Path
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		actualBody = string(body)

		response := `{ 
    		"seq": 2, 
    		"error": 0,
    		"data": {
        		"switch": "on",
        		"startup": "off",
        		"pulse": "off",
        		"pulseWidth": 500,
        		"ssid": "eWeLink",
        		"otaUnlock": false,
        		"fwVersion": "3.5.0",
        		"deviceid": "100000140e",
        		"bssid": "ec:17:2f:3d:15:e",
        		"signalStrength": -25
			}
 		}`
		rw.Write([]byte(response))
	}))

	srv.Listener.Close()
	srv.Listener = l

	srv.Start()
	defer srv.Close()

	actualState, err := testAPI.GetState()

	assert.NoError(t, err)
	assert.Equal(t, expectedState, actualState)
	assert.Equal(t, expectedURLPath, actualURLPath)
	assert.Equal(t, expectedBody, actualBody)
}

func Test_GetState_Off(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	l, err := net.Listen("tcp", "127.0.0.1:8081")
	assert.NoError(t, err)

	expectedState := "off"

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		response := `{ 
    		"seq": 2, 
    		"error": 0,
    		"data": {
        		"switch": "off",
        		"startup": "off",
        		"pulse": "off",
        		"pulseWidth": 500,
        		"ssid": "eWeLink",
        		"otaUnlock": false,
        		"fwVersion": "3.5.0",
        		"deviceid": "100000140e",
        		"bssid": "ec:17:2f:3d:15:e",
        		"signalStrength": -25
			}
 		}`
		rw.Write([]byte(response))
	}))

	srv.Listener.Close()
	srv.Listener = l

	srv.Start()
	defer srv.Close()

	actualState, err := testAPI.GetState()

	assert.NoError(t, err)
	assert.Equal(t, expectedState, actualState)
}

func Test_GetState_ErrorInResponse(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	l, err := net.Listen("tcp", "127.0.0.1:8081")
	assert.NoError(t, err)

	expectedState := ""

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		response := `{ 
    		"seq": 2, 
    		"error": 1,
    		"data": {
        		"switch": "off",
        		"startup": "off",
        		"pulse": "off",
        		"pulseWidth": 500,
        		"ssid": "eWeLink",
        		"otaUnlock": false,
        		"fwVersion": "3.5.0",
        		"deviceid": "100000140e",
        		"bssid": "ec:17:2f:3d:15:e",
        		"signalStrength": -25
			}
 		}`
		rw.Write([]byte(response))
	}))

	srv.Listener.Close()
	srv.Listener = l

	srv.Start()
	defer srv.Close()

	actualState, err := testAPI.GetState()

	assert.Error(t, err)
	assert.Equal(t, expectedState, actualState)
}

func Test_GetState_ErrorNoServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	expectedState := ""
	actualState, err := testAPI.GetState()

	assert.Error(t, err)
	assert.Equal(t, expectedState, actualState)
}

func Test_GetState_ErrorGettingIPAddr(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("", errors.New("dummy error")).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	expectedState := ""
	actualState, err := testAPI.GetState()

	assert.Error(t, err)
	assert.Equal(t, expectedState, actualState)
}

func Test_TurnOn(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	l, err := net.Listen("tcp", "127.0.0.1:8081")
	assert.NoError(t, err)

	expectedURLPath := "/zeroconf/switch"
	actualURLPath := ""
	expectedBody := `{ 
		"deviceid": "", 
		"data": {
			"switch": "on" 
		} 
	 }`
	actualBody := ""

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		actualURLPath = r.URL.Path
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		actualBody = string(body)

		response := `{ 
			"error": 0
		 }`
		rw.Write([]byte(response))
	}))

	srv.Listener.Close()
	srv.Listener = l

	srv.Start()
	defer srv.Close()

	err = testAPI.TurnOn()

	assert.NoError(t, err)
	assert.Equal(t, expectedURLPath, actualURLPath)
	assert.Equal(t, expectedBody, actualBody)
}

func Test_TurnOn_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	l, err := net.Listen("tcp", "127.0.0.1:8081")
	assert.NoError(t, err)

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		response := `{ 
			"error": 1
		 }`
		rw.Write([]byte(response))
	}))

	srv.Listener.Close()
	srv.Listener = l

	srv.Start()
	defer srv.Close()

	err = testAPI.TurnOn()

	assert.Error(t, err)
}

func Test_TurnOn_ErrorNoServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	err := testAPI.TurnOn()

	assert.Error(t, err)
}

func Test_TurnOff(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	l, err := net.Listen("tcp", "127.0.0.1:8081")
	assert.NoError(t, err)

	expectedURLPath := "/zeroconf/switch"
	actualURLPath := ""
	expectedBody := `{ 
		"deviceid": "", 
		"data": {
			"switch": "off" 
		} 
	 }`
	actualBody := ""

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		actualURLPath = r.URL.Path
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		actualBody = string(body)

		response := `{ 
			"error": 0
		 }`
		rw.Write([]byte(response))
	}))

	srv.Listener.Close()
	srv.Listener = l

	srv.Start()
	defer srv.Close()

	err = testAPI.TurnOff()

	assert.NoError(t, err)
	assert.Equal(t, expectedURLPath, actualURLPath)
	assert.Equal(t, expectedBody, actualBody)
}

func Test_TurnOff_Error(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	l, err := net.Listen("tcp", "127.0.0.1:8081")
	assert.NoError(t, err)

	srv := httptest.NewUnstartedServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		response := `{ 
			"error": 1
		 }`
		rw.Write([]byte(response))
	}))

	srv.Listener.Close()
	srv.Listener = l

	srv.Start()
	defer srv.Close()

	err = testAPI.TurnOff()

	assert.Error(t, err)
}

func Test_TurnOff_ErrorNoServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	osSvcs := mockosservicesprovider.NewMockOSServicesProvider(mockCtrl)
	osSvcs.EXPECT().GetIPFromMAC("dummy_mac_address").Return("127.0.0.1", nil).Times(1)

	testAPI := New(osSvcs, "dummy_mac_address")

	err := testAPI.TurnOff()

	assert.Error(t, err)
}
