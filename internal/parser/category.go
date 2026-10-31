package parser

import "strings"

const (
	CatFood          = "Makanan & Minuman"
	CatTransport     = "Transportasi"
	CatKos           = "Kos & Utilitas"
	CatPhone         = "Pulsa & Internet"
	CatKuliah        = "Kuliah"
	CatShopping      = "Belanja"
	CatEntertainment = "Hiburan"
	CatHealth        = "Kesehatan"
	CatSocial        = "Sosial & Hadiah"
	CatOther         = "Lainnya"

	CatAllowance   = "Uang Saku"
	CatSalary      = "Gaji/Freelance"
	CatScholarship = "Beasiswa"
)

var ExpenseCategories = []string{
	CatFood, CatTransport, CatKos, CatPhone, CatKuliah, CatShopping,
	CatEntertainment, CatHealth, CatSocial, CatOther,
}

var IncomeCategories = []string{CatAllowance, CatSalary, CatScholarship, CatOther}

var expenseKeywords = map[string][]string{
	CatFood: {"makan", "bakso", "nasi", "mie", "mi", "ayam", "soto", "sate", "es teh", "teh", "kopi",
		"kopken", "minum", "jajan", "snack", "gorengan", "martabak", "seblak", "galon", "warteg",
		"sarapan", "boba", "roti", "indomie", "geprek", "pecel", "bubur"},
	CatTransport: {"bensin", "ojek", "ojol", "gojek", "grab", "maxim", "parkir", "bus", "krl",
		"kereta", "angkot", "transjakarta", "tol", "taksi", "pertalite", "pertamax"},
	CatKos:           {"kos", "kost", "listrik", "token", "air", "wifi", "laundry", "sewa"},
	CatPhone:         {"pulsa", "kuota", "paket data", "internet"},
	CatKuliah:        {"print", "fotokopi", "fotocopy", "buku", "ukt", "spp", "modul", "praktikum", "jilid", "tugas"},
	CatShopping:      {"baju", "sepatu", "celana", "shopee", "tokopedia", "sabun", "sampo", "skincare", "belanja", "indomaret", "alfamart"},
	CatEntertainment: {"nonton", "bioskop", "netflix", "spotify", "game", "topup", "top up", "karaoke", "konser"},
	CatHealth:        {"obat", "dokter", "apotek", "vitamin", "klinik", "masker"},
	CatSocial:        {"kado", "hadiah", "sumbangan", "infaq", "sedekah", "traktir", "kondangan", "patungan"},
}

var incomeKeywords = map[string][]string{
	CatAllowance:   {"ortu", "mama", "papa", "ibu", "ayah", "bapak", "uang saku", "kiriman"},
	CatSalary:      {"gaji", "gajian", "freelance", "fee", "honor", "proyek", "project"},
	CatScholarship: {"beasiswa"},
}

// incomeWords mark a message as income.
var incomeWords = []string{"dapet", "dapat", "gajian", "transferan", "dikasih", "cair", "beasiswa", "terima", "masuk"}

func IsIncome(text string) bool { return containsAny(text, incomeWords) }

// Categorize returns a default category and whether a keyword was recognised.
func Categorize(text string, t TransactionType) (string, bool) {
	keywords, order := expenseKeywords, ExpenseCategories
	if t == TypeIncome {
		keywords, order = incomeKeywords, IncomeCategories
	}
	for _, cat := range order {
		if containsAny(text, keywords[cat]) {
			return cat, true
		}
	}
	return CatOther, false
}

// containsAny matches whole words (or phrases), case-insensitively.
func containsAny(text string, words []string) bool {
	padded := " " + strings.Join(strings.FieldsFunc(strings.ToLower(text), isSeparator), " ") + " "
	for _, w := range words {
		if strings.Contains(padded, " "+w+" ") {
			return true
		}
	}
	return false
}

func isSeparator(r rune) bool {
	return r == ' ' || r == ',' || r == '.' || r == '?' || r == '!' || r == '\n' || r == '\t'
}
