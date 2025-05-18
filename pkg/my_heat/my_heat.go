package my_heat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func MyHeatGetDevices(login, key string) ([]*Device, error) {
	timeout := os.Getenv("DEFAULT_REQUEST_TIMEOUT")
	_timeout, err := strconv.Atoi(timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout for request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(_timeout)*time.Second)
	defer cancel()

	body := MyHeatRequestBody{
		Action: "getDevices",
		Login:  login,
		Key:    key,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error while marshal request body%v", err)
	}

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "POST", os.Getenv("MYHEAT_ENDPOINT"), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error while create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while do request: %v", err)
	}

	var res *MyHeatResponseBody
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status-code %v", resp.StatusCode)

	} else {
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("error while decode response body%v", err)
		}

		switch res.Err {
		case 0:
			return res.Data.Devices, nil
		default:
			return nil, fmt.Errorf("getting devices list filed: error code %d", res.Err)
		}
	}
}

func MyHeatGetDeviceInfo(login, key string, deviceID int) (*Data, error) {
	timeout := os.Getenv("DEFAULT_REQUEST_TIMEOUT")
	_timeout, err := strconv.Atoi(timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout for request: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(_timeout)*time.Second)
	defer cancel()

	body := MyHeatRequestBody{
		Action:   "getDeviceInfo",
		Login:    login,
		Key:      key,
		DeviceID: deviceID,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("error while marshal request body%v", err)
	}

	client := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "POST", os.Getenv("MYHEAT_ENDPOINT"), bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error while create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error while do request: %v", err)
	}

	var res *MyHeatResponseBody
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrong status-code %v", resp.StatusCode)

	} else {
		err = json.NewDecoder(resp.Body).Decode(&res)
		if err != nil {
			return nil, fmt.Errorf("error while decode response body%v", err)
		}

		switch res.Err {
		case 0:
			return res.Data, nil
		default:
			return nil, fmt.Errorf("getting devices list filed: error code %d", res.Err)
		}
	}
}
