package qzone

type UploadPicResult struct {
	Data struct {
		Pre          string `json:"pre"`
		URL          string `json:"url"`
		Lloc         string `json:"lloc"`
		Sloc         string `json:"sloc"`
		Type         int    `json:"type"`
		Width        int    `json:"width"`
		Height       int    `json:"height"`
		Albumid      string `json:"albumid"`
		Totalpic     int    `json:"totalpic"`
		Limitpic     int    `json:"limitpic"`
		OriginURL    string `json:"origin_url"`
		OriginUUID   string `json:"origin_uuid"`
		OriginWidth  int    `json:"origin_width"`
		OriginHeight int    `json:"origin_height"`
		Contentlen   int    `json:"contentlen"`
	} `json:"data"`
	Ret int `json:"ret"`
}
type ShuoshuoList struct {
	Code    int    `json:"code"`
	Subcode int    `json:"subcode"`
	Message string `json:"message"`
	Default int    `json:"default"`
	Data    struct {
		Main struct {
			Attach               string        `json:"attach"`
			Searchtype           string        `json:"searchtype"`
			Hasmorefeeds         bool          `json:"hasMoreFeeds"`
			Daylist              string        `json:"daylist"`
			Uinlist              string        `json:"uinlist"`
			Error                string        `json:"error"`
			Hotkey               string        `json:"hotkey"`
			Icgroupdata          []interface{} `json:"icGroupData"`
			HostLevel            string        `json:"host_level"`
			FriendLevel          string        `json:"friend_level"`
			Lastaccesstime       string        `json:"lastaccesstime"`
			Lastaccessrelatetime string        `json:"lastAccessRelateTime"`
			Begintime            string        `json:"begintime"`
			Endtime              string        `json:"endtime"`
			Dayspac              string        `json:"dayspac"`
			Hidednamelist        []interface{} `json:"hidedNameList"`
			Aisortbegintime      string        `json:"aisortBeginTime"`
			Aisortendtime        string        `json:"aisortEndTime"`
			Aisortoffset         string        `json:"aisortOffset"`
			Aisortnexttime       string        `json:"aisortNextTime"`
			OwnerBitmap          string        `json:"owner_bitmap"`
			Pagenum              string        `json:"pagenum"`
			Externparam          string        `json:"externparam"`
		} `json:"main"`
		Data []map[string]interface{} `json:"data"`
	} `json:"data"`
}
type SendShuoShuoResult struct {
	Attach     string `json:"attach"`
	Code       int    `json:"code"`
	Feedinfo   string `json:"feedinfo"`
	Message    string `json:"message"`
	Needverify int    `json:"needVerify"`
	Now        int    `json:"now"`
	Republish  int    `json:"republish"`
	Secret     int    `json:"secret"`
	Subcode    int    `json:"subcode"`
	Tid        string `json:"tid"`
	Vote       string `json:"vote"`
}
