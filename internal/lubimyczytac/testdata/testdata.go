package testdata

import (
	_ "embed"

	"github.com/Kichiyaki/lubimyczytacrss/internal/lubimyczytac"
)

//go:embed remigiusz_mroz.html
var remigiuszMrozHTML []byte

//go:embed john_flanagan.html
var johnFlanaganHTML []byte

type AuthorWithHTML struct {
	lubimyczytac.Author
	HTML []byte
}

var Authors = []AuthorWithHTML{
	{
		Author: lubimyczytac.Author{
			ID:               "82094",
			Name:             "Remigiusz Mróz",
			ShortDescription: "Polski pisarz, autor powieści kryminalnych oraz cyklu publicystycznego „Kurs pisania”.Ukończył z wyróżnieniem Akademię Leona Koźmińskiego w Warszawie, gd...",
			URL:              "https://lubimyczytac.pl/autor/82094/remigiusz-mroz",
			Books: []lubimyczytac.Book{
				{Title: "Obrazy z przeszłości", URL: "https://lubimyczytac.pl/ksiazka/5016998/obrazy-z-przeszlosci"},
				{Title: "Skazanie", URL: "https://lubimyczytac.pl/ksiazka/5009453/skazanie"},
				{Title: "Behawiorysta", URL: "https://lubimyczytac.pl/ksiazka/5006528/behawiorysta"},
				{Title: "Projekt Riese", URL: "https://lubimyczytac.pl/ksiazka/4998407/projekt-riese"},
				{Title: "Immunitet", URL: "https://lubimyczytac.pl/ksiazka/4990168/immunitet"},
				{Title: "Przepaść", URL: "https://lubimyczytac.pl/ksiazka/4988766/przepasc"},
				{Title: "Egzekucja", URL: "https://lubimyczytac.pl/ksiazka/4983192/egzekucja"},
				{Title: "Wybaczam ci", URL: "https://lubimyczytac.pl/ksiazka/4975543/wybaczam-ci"},
				{Title: "Inwigilacja", URL: "https://lubimyczytac.pl/ksiazka/4968996/inwigilacja"},
				{Title: "Ekstremista", URL: "https://lubimyczytac.pl/ksiazka/4968712/ekstremista"},
				{Title: "Afekt", URL: "https://lubimyczytac.pl/ksiazka/4962195/afekt"},
				{Title: "Głębia osobliwości cz. 2", URL: "https://lubimyczytac.pl/ksiazka/5009924/glebia-osobliwosci-cz-2"},
				{Title: "Szepty spoza nicości", URL: "https://lubimyczytac.pl/ksiazka/4955433/szepty-spoza-nicosci"},
				{Title: "W cieniu prawa", URL: "https://lubimyczytac.pl/ksiazka/4947287/w-cieniu-prawa"},
				{Title: "Księgarenka przy ulicy Wiśniowej", URL: "https://lubimyczytac.pl/ksiazka/4944645/ksiegarenka-przy-ulicy-wisniowej"},
				{Title: "Rewizja", URL: "https://lubimyczytac.pl/ksiazka/4943912/rewizja"},
				{Title: "Halny", URL: "https://lubimyczytac.pl/ksiazka/4943948/halny"},
				{Title: "Precedens", URL: "https://lubimyczytac.pl/ksiazka/4939105/precedens"},
				{Title: "Osiedle RZNiW", URL: "https://lubimyczytac.pl/ksiazka/4932748/osiedle-rzniw"},
				{Title: "Lot 202", URL: "https://lubimyczytac.pl/ksiazka/4926011/lot-202"},
			},
		},
		HTML: remigiuszMrozHTML,
	},
	{
		Author: lubimyczytac.Author{
			ID:               "19013",
			Name:             "John Flanagan",
			ShortDescription: "John Flanagan Urodzony i wychowany w Sydney w Australii, John Flanagan od dzieciństwa marzył o tym, by zostać pisarzem. Nie było łatwo. Pracował w agencji r...",
			URL:              "https://lubimyczytac.pl/autor/19013/john-flanagan",
			Books: []lubimyczytac.Book{
				{Title: "Morska pogoń", URL: "https://lubimyczytac.pl/ksiazka/5016250/morska-pogon"},
				{Title: "Ruiny Gorlanu", URL: "https://lubimyczytac.pl/ksiazka/4989803/ruiny-gorlanu"},
				{Title: "Ziemia skuta lodem", URL: "https://lubimyczytac.pl/ksiazka/4989807/ziemia-skuta-lodem"},
				{Title: "Bitwa o Skandię", URL: "https://lubimyczytac.pl/ksiazka/4989809/bitwa-o-skandie"},
				{Title: "Ucieczka z zamku Falaise", URL: "https://lubimyczytac.pl/ksiazka/4984975/ucieczka-z-zamku-falaise"},
				{Title: "Płonący most", URL: "https://lubimyczytac.pl/ksiazka/4989806/plonacy-most"},
				{Title: "Zaginiony książę", URL: "https://lubimyczytac.pl/ksiazka/4935503/zaginiony-ksiaze"},
				{Title: "Powrót Temudżeinów", URL: "https://lubimyczytac.pl/ksiazka/4897979/powrot-temudzeinow"},
				{Title: "Pojedynek w Araluenie", URL: "https://lubimyczytac.pl/ksiazka/4861058/pojedynek-w-araluenie"},
				{Title: "Klan Czerwonego Lisa", URL: "https://lubimyczytac.pl/ksiazka/4851290/klan-czerwonego-lisa"},
				{Title: "Kaldera", URL: "https://lubimyczytac.pl/ksiazka/4811916/kaldera"},
				{Title: "Bitwa na Wrzosowiskach", URL: "https://lubimyczytac.pl/ksiazka/3874675/bitwa-na-wrzosowiskach"},
				{Title: "Nieznany ląd", URL: "https://lubimyczytac.pl/ksiazka/303434/nieznany-lad"},
				{Title: "Turniej w Gorlanie", URL: "https://lubimyczytac.pl/ksiazka/266058/turniej-w-gorlanie"},
				{Title: "Góra Skorpiona", URL: "https://lubimyczytac.pl/ksiazka/230915/gora-skorpiona"},
				{Title: "Niewolnicy z Socorro", URL: "https://lubimyczytac.pl/ksiazka/220609/niewolnicy-z-socorro"},
				{Title: "Królewski zwiadowca", URL: "https://lubimyczytac.pl/ksiazka/192775/krolewski-zwiadowca"},
				{Title: "Pościg", URL: "https://lubimyczytac.pl/ksiazka/167327/poscig"},
				{Title: "Najeźdźcy", URL: "https://lubimyczytac.pl/ksiazka/144610/najezdzcy"},
				{Title: "Zaginione historie", URL: "https://lubimyczytac.pl/ksiazka/131440/zaginione-historie"},
			},
		},
		HTML: johnFlanaganHTML,
	},
}
