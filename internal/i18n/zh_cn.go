package i18n

// zhCN is the Simplified Chinese catalog. The key is the English text, and
// "text|context" for TC.
var zhCN = map[string]string{
	// pages
	"Read":      "阅读",
	"Devotions": "灵修",
	"Pray":      "祷告",
	"Stories":   "故事",
	"Seasons":   "节期",
	"Music":     "音乐",
	"Device":    "设备",

	// Ago and friends
	"just now":       "刚刚",
	"a minute ago":   "一分钟前",
	"%d minutes ago": "%d 分钟前",
	"an hour ago":    "一小时前",
	"%d hours ago":   "%d 小时前",
	"yesterday":      "昨天",
	"%d days ago":    "%d 天前",

	// shared words
	"Settings":             "设置",
	"Language":             "语言",
	"System":               "跟随系统",
	"Save":                 "保存",
	"Save|note":            "保存",
	"Save|prayers":         "记下",
	"Cancel":               "取消",
	"Done":                 "完成",
	"Search":               "搜索",
	"Next":                 "下一个",
	"Back":                 "返回",
	"Volume":               "音量",
	"Retry":                "重试",
	"Not paired":           "未配对",
	"Verse of the Day":     "每日金句",
	"Morning & Evening":    "晨昏灵修",
	"Reading Plans":        "读经计划",
	"Prayer List":          "祷告清单",
	"This week's readings": "本周经课",
	"Play":                 "播放",
	"Pause":                "暂停",
	"Listen":               "朗读",

	// the device page and its pairing guide
	"No device is paired. Pair one to hear the Bible read aloud.": "尚未配对设备。配对一台，即可听它朗读圣经。",
	"Paired with %s.":     "已与 %s 配对。",
	"Pair a device":       "配对设备",
	"Pair a device first": "请先配对设备",
	"Not now":             "暂不",
	"Pair":                "配对",
	"the device":          "设备",
	"Character":           "形象",
	"Action group":        "动作组",
	"The device plays its own motions as it speaks. How strong they read is yours to set.": "设备说话时会自行播放自带动作，动作幅度由你设定。",
	"Make your Codie as Bibel": "让你的 Codie 成为圣经伙伴",

	// the pairing guide
	"Pair your codie device, and it reads the Bible aloud: the stories, the psalms, the verses of the day. It speaks as Jesus, or as a monk for the quieter rooms.": "配对你的 codie 设备，让它朗读圣经：故事、诗篇、每日金句。它以耶稣的形象说话，也可换成修士，适合安静的场合。",
	"Turn the device on, and join the same Wi-Fi as this computer.": "打开设备，并让它与这台电脑连到同一个 Wi-Fi。",
	"The device's address": "设备地址",
	"Type the address the device's screen shows, as 192.168.1.10:8081.": "输入设备屏幕上显示的地址，如 192.168.1.10:8081。",
	"Device address":                                  "设备地址",
	"Asking the device…":                              "正在询问设备……",
	"Type the device's address first.":                "请先输入设备地址。",
	"The pairing code":                                "配对码",
	"Pairing code":                                    "配对码",
	"Type the six figures the device's screen shows.": "输入设备屏幕上显示的六位数字。",
	"The code is six figures.":                        "配对码是六位数字。",
	"Pairing…":                                        "正在配对……",
	"Who speaks?":                                     "谁来说话？",
	"Dressing the character…":                         "正在装扮形象……",
	"%s speaks on the device now.":                    "%s 已在设备上说话了。",
	"No model file for %s is in the app's data folder, so the device speaks with its own face. Put %s there to change that.": "应用数据目录里没有 %s 的模型文件，设备将用它自带的面孔说话。把 %s 放进该目录即可更换。",

	// page placeholders, until each page arrives
	"The Bible, book by book, chapter by chapter.":                  "圣经,一卷一卷,一章一章。",
	"The verse of the day, Morning and Evening, and reading plans.": "每日金句、晨昏灵修与读经计划。",
	"A guided prayer, and the prayers you keep.":                    "引导式祷告,和你记下的祷告。",
	"The Bible's stories, told aloud.":                              "圣经故事,讲给你听。",
	"The church year, its feasts, and this week's readings.":        "教会年历、节期与本周经课。",
	"Hymns to play quietly in the background.":                      "轻声播放的赞美诗。",
	"Pair a codie device to read the Bible aloud.":                  "配对 codie 设备,朗读圣经。",
	"Language, appearance, and the reader.":                         "语言、外观与阅读器。",
}
