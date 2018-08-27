package document

const DocTypeSprecher = "sprecher"

type Base struct {
	ID        string
	DocType   string
	Published bool
}

func (p *Base) StoredFields() (fields []string) {
	fields = append(fields, []string{"Published"}...)
	return fields
}

type PartCount struct {
	Base
	AnzahlRollen uint64
}

func (p *PartCount) StoredFields() (fields []string) {
	fields = append(fields, []string{"AnzahlRollen"}...)
	return append(fields, p.Base.StoredFields()...)
}

type Person struct {
	PartCount
	Anrede         string
	Nachname       string
	Vorname        string
	Zusatz         string
	PseudoVorname  string
	PseudoNachname string
	Name           string
	PseudoName     string
}

func (p *Person) StoredFields() (fields []string) {
	fields = append(fields, []string{"Vorname", "Nachname", "Anrede", "Zusatz", "PseudoNachname", "PseudoVorname"}...)
	return append(fields, p.PartCount.StoredFields()...)
}

type Sprecher struct {
	Person
	Beschreibung string
	Geburtsort   string
	Todesort     string
}

func (s *Sprecher) StoredFields() (fields []string) {
	return append(fields, s.Person.StoredFields()...)
}

func (Sprecher) Type() string {
	return "sprecher"
}
