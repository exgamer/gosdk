package constants

const LangRu string = "ru"
const LangKZ string = "kk"
const LangEN string = "en"

func GetLanguages() []string {
	return []string{
		LangRu,
		LangKZ,
		LangEN,
	}
}
