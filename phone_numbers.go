package kazooapi

import (
	"context"
	"encoding/json"
	"io/ioutil"
)

type PhoneNumbersAPIService service

type (
	PhoneNumber struct {
		ID       string   `json:"id"`
		State    string   `json:"state"`
		Features []string `json:"features"`
		AssignedTo string `json:"assigned_to"`
		Created Timestamp `json:"created"`
		Updated Timestamp `json:"updated"`
		readOnly struct {
			state    string   `json:"state,omitempty"`
			created string `json:"created,omitempty"`
			modified string `json:"modified,omitempty"`
			features []string `json:"features,omitempty"`
		} `json:"_read_only,omitempty"`
	}

	AccountPhoneNumbersResponse struct {
		CascadeQuantity int64 `json:"cascade_quantity"`
		Numbers map[string]PhoneNumber `json:"numbers"`
	}
)

func (api *PhoneNumbersAPIService) CreatePhoneNumber(ctx context.Context, acc string, num string) (number *PhoneNumber, err error) {
	var response struct {
		Data PhoneNumber `json:"data"`
		ResponseEnvelope
	}

	if num == "" {
		return nil, reportError("number is required field")
	}

	params := Request{
		CTX:    ctx,
		Method: "PUT",
		Path:   api.client.cfg.BasePath + "/accounts/" + acc + "/phone_numbers/" + num,
	}

/*	reqBody := RequestEnvelope{
		Data: input,
	}

	jsonString, err := json.Marshal(reqBody)
	if err != nil {
		return nil, reportError("can't marshall body for request")
	}

	body, err := setBody(jsonString, "json")
	if err != nil {
		return nil, reportError("can't prepare body for the request")
	}

	params.PostBody = body*/

	req, err := api.client.prepareRequest(&params)
	if err != nil {
		return nil, reportError("Can't prepare a request %s", err)
	}

	resp, err := api.client.callAPI(ctx, req)
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, reportError("Status: %v, Body: %s", resp.Status, bodyBytes)
	}

	decoder := json.NewDecoder(resp.Body)

	decErr := decoder.Decode(&response)
	if decErr != nil {
		return nil, reportError("Can't decode response: %v", decErr)
	}

	number = &response.Data

	return number, nil
}

func (api *PhoneNumbersAPIService) ListPhoneNumbers(ctx context.Context, acc string, disablePagination bool) (numbers []PhoneNumber, err error) {
	var response struct {
		Data AccountPhoneNumbersResponse `json:"data"`
		ResponseEnvelope
	}

	params := Request{
		CTX:    ctx,
		Method: "GET",
		Path:   api.client.cfg.BasePath + "/accounts/" + acc + "/phone_numbers",
	}

	if disablePagination {
		params.QueryParams = map[string][]string{"paginate": []string{"false"}}
	}

	req, err := api.client.prepareRequest(&params)
	if err != nil {
		return nil, reportError("Can't prepare a request %s", err)
	}

	resp, err := api.client.callAPI(ctx, req)
	if err != nil || resp == nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return nil, reportError("Status: %v, Body: %s", resp.Status, bodyBytes)
	}

	decoder := json.NewDecoder(resp.Body)

	decErr := decoder.Decode(&response)
	if decErr != nil {
		return nil, reportError("Can't decode response: %v", decErr)
	}

	for key, number := range response.Data.Numbers {
		number.ID = key
		numbers = append(numbers, number)
	}

	return numbers, nil
}