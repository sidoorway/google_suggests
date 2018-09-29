package gsuggest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
)

/*
["поисковые подсказки ",[["поисковые подсказки\u003cb\u003e google не работают\u003c\/b\u003e",0],["поисковые подсказки\u003cb\u003e гугл\u003c\/b\u003e",0,[67]],["поисковые подсказки\u003cb\u003e яндексе\u003c\/b\u003e",0],["поисковые подсказки\u003cb\u003e seo\u003c\/b\u003e",0,[22,30]],["поисковые подсказки\u003cb\u003e накрутка\u003c\/b\u003e",0,[22,30]],["поисковые подсказки\u003cb\u003e bitrix\u003c\/b\u003e",0,[22,30]],["поисковые подсказки\u003cb\u003e отключить\u003c\/b\u003e",0,[22,30]],["\u003cb\u003eбитрикс \u003c\/b\u003eпоисковые подсказки",0,[22,30]],["\u003cb\u003evivaldi \u003c\/b\u003eпоисковые подсказки",0,[22,30]],["\u003cb\u003eвключить \u003c\/b\u003eпоисковые подсказки",0,[22,30]]],{"i":"поисковые подсказки","j":"b","q":"z-VGWq9jHq-JbZZA_wHnFj_Dr6U","t":{"bpc":false,"tlw":false}}]
*/

// https://www.google.ru/complete/search?
// client=psy-ab&
// hl=ru&
// gs_rn=64&
// gs_ri=psy-ab&
// gs_mss=%D0%BF%D0%BE%D0%B8%D1%81%D0%BA%D0%BE%D0%B2%D1%8B%D0%B5%20%D0%BF%D0%BE%D0%B4%D1%81%D0%BA%D0%B0%D0%B7%D0%BA%D0%B8&
// newwindow=1&
// ei=VnGuW_-NN8i7sQHwwLzYAQ&
// pq=%D0%BF%D0%BE%D0%B8%D1%81%D0%BA%D0%BE%D0%B2%D1%8B%D0%B5%20%D0%BF%D0%BE%D0%B4%D1%81%D0%BA%D0%B0%D0%B7%D0%BA%D0%B8%20github&
// cp=20&
// gs_id=b&
// q=%D0%BF%D0%BE%D0%B8%D1%81%D0%BA%D0%BE%D0%B2%D1%8B%D0%B5%20%D0%BF%D0%BE%D0%B4%D1%81%D0%BA%D0%B0%D0%B7%D0%BA%D0%B8%20&
// xhr=t

// Get - retrieve google suggestions
func Get(keyword, lang, country string) ([]string, error) {

	link, err := url.Parse("https://www.google.ru/complete/search")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("client", "psy_ab")
	params.Add("hl", lang)
	params.Add("gl", country)
	params.Add("gs_rn", "64")
	params.Add("gs_ri", "psy-ab")
	params.Add("gs_mss", keyword)
	params.Add("newwindow", "1")
	params.Add("ei", "VnGuW_-NN8i7sQHwwLzYAQ")
	params.Add("pq", keyword)
	params.Add("cp", "20")
	params.Add("gs_id", "b")
	params.Add("q", keyword)
	params.Add("xhr", "t")
	params.Add("ie", "utf-8")
	params.Add("oe", "utf-8")

	link.RawQuery = params.Encode()

	resp, err := http.Get(link.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	gresp := make([]interface{}, 0)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(buf).Decode(&gresp)
	if err != nil {
		return nil, err
	}

	suggests := make([]string, 0)
	i := 0
	for {
		if s1, ok := gresp[1].([]interface{}); ok {
			if i < len(s1) {
				if s2, ok := s1[i].([]interface{}); ok {
					if s3, ok := s2[0].(string); ok {
						suggests = append(suggests, s3)
						i++
						continue
					}
				}
			}
		}
		break
	}

	return suggests, nil
}
