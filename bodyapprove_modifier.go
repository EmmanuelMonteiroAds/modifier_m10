package modifiermcabodyapprove

import (	
	"net/http"	
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"bytes"
	"strings"

	"github.com/google/martian/parse"
)

func init() {
	parse.Register("bodymcaapprove.Modifier", modifierFromJSON)
}

type XmlModifier struct {
	contentType string
}

type XmlModifierJSON struct {
	ContentType string               `json:"contentType"`
	Scope       []parse.ModifierType `json:"scope"`
}

type Request struct {
    XMLName   xml.Name `xml:"Request" json:"-"`
    SafrapayId  string   `xml:"safrapayId" json:"safrapayId"`
	PartiallyApproved  bool   `xml:"partiallyApproved" json:"partiallyApproved"`
	RiskRating  string   `xml:"riskRating" json:"riskRating"`
	CreditScore  string   `xml:"creditScore" json:"creditScore"`
	AdverseAction   *AdverseAction `xml:"adverseAction" json:"adverseAction"`
    McaOfferOne   *McaOfferOne `xml:"mcaOfferOne" json:"mcaOfferOne"`
	McaOfferTwo   *McaOfferTwo `xml:"mcaOfferTwo" json:"mcaOfferTwo"`
	ApplicationId string   `xml:"applicationId" json:"applicationId"`
	ProductDesc  string   `xml:"productDesc" json:"productDesc"`
}

type AdverseAction struct {
    Code  string `xml:"code" json:"code"`
    Source string `xml:"source" json:"source"`
    Description string `xml:"description" json:"description"`
}

type McaOfferOne struct {
    Key  int `xml:"offerKey" json:"key"`
    WithHold float32 `xml:"withHold" json:"withHold"`
    AdvanceAmount float32 `xml:"advanceAmount" json:"advanceAmount"`
	PaybackAmount float32 `xml:"paybackAmount" json:"paybackAmount"`
	Duration  int `xml:"duration" json:"duration"`
}

type McaOfferTwo struct {
    Key  int `xml:"offerKey" json:"key"`
    WithHold float32 `xml:"withHold" json:"withHold"`
    AdvanceAmount float32 `xml:"advanceAmount" json:"advanceAmount"`
	PaybackAmount float32 `xml:"paybackAmount" json:"paybackAmount"`
	Duration  int `xml:"duration" json:"duration"`
}

func (m *XmlModifier) ModifyRequest(req *http.Request) error {
	
	req.Header.Set("Content-Type", m.contentType)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
    var data Request
    xml.Unmarshal([]byte(body), &data)
    jsonData, _ := json.Marshal(data)

	req.ContentLength = int64(len(jsonData))
	
	req.Body = ioutil.NopCloser(bytes.NewReader(jsonData))
  
  	req.URL.Path = strings.Replace(req.URL.Path, "safrapayId", data.SafrapayId, 1);

	return nil
}

func (m *XmlModifier) ModifyResponse(res *http.Response) error {
	
	if(res.StatusCode != 200){
		res.StatusCode = 200
	}	

	return nil
}

func XmlNewModifier(contentType string) *XmlModifier {
	return &XmlModifier{
		contentType:  contentType,
	}
}

func modifierFromJSON(b []byte) (*parse.Result, error) {
	msg := &XmlModifierJSON{}

	if err := json.Unmarshal(b, msg); err != nil {
		return nil, err
	}

	return parse.NewResult(XmlNewModifier(msg.ContentType), msg.Scope)
}