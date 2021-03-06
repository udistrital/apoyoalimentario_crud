package models

import (
	"apoyoalimentario_CRUD_API/db"
	"apoyoalimentario_CRUD_API/utility"
	"strings"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Economic Struct for save information economic of student(s)
type Economic struct {
	ID              bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Idc             bson.ObjectId `json:"id" bson:"id"`
	Estrato         string        `json:"estrato" bson:"estrato"`
	Ingresos        int           `json:"ingresos" bson:"ingresos"`
	SostePropia     string        `json:"sostenibilidadpropia" bson:"sostenibilidadpropia"`
	SosteHogar      string        `json:"sostenibilidadhogar" bson:"sostenibilidadhogar"`
	Nucleofam       string        `json:"nucleofamiliar" bson:"nucleofamiliar"`
	PersACargo      string        `json:"personasacargo" bson:"personasacargo"`
	EmpleadArriendo string        `json:"empleadoroarriendo" bson:"empleadoroarriendo"`
	ProvBogota      string        `json:"provienefuerabogota" bson:"provienefuerabogota"`
	Ciudad          string        `json:"ciudad" bson:"ciudad"`
	PobEspecial     string        `json:"poblacionespecial" bson:"poblacionespecial"`
	Discapacidad    string        `json:"discapacidad" bson:"discapacidad"`
	PatAlimenticia  string        `json:"patologiaalimenticia" bson:"patologiaalimenticia"`
	SerPiloPaga     string        `json:"serpilopaga" bson:"serpilopaga"`
	Sisben          string        `json:"sisben" bson:"sisben"`
	Periodo         int           `json:"periodo" bson:"periodo"`
	SemestreIns     int           `json:"semestre" bson:"semestre"`
	Matricula       int           `json:"matricula" bson:"matricula"`
	EstadoProg      int           `json:"estadoprograma" bson:"estadoprograma"`
	TipoSubsidio    string        `json:"tiposubsidio" bson:"tiposubsidio"`
	Tipoapoyo       string        `json:"tipoapoyo" bson:"tipoapoyo"`
	Mensaje         string        `json:"mensaje" bson:"mensaje"`
	Telefono        string        `json:"telefono" bson:"telefono"`
	Correo          string        `json:"correo" bson:"correo"`
	Antiguedad      string        `json:"antiguedad" bson:"antiguedad"`
	Verificadopor   string        `json:"verificadopor" bson:"verificadopor"`
}

//GetInformationEconomic - Optine la informacion economica de un estudiante en el actual semestre
//Param session		IN   "sesion de base de datos"
//Param code		IN   "Codigo del estudiante a consultar"
//Param InfoEcono	OUT   "modelo de informacion economica consultado"
//Param errP		OUT   "error si es que existe"
func GetInformationEconomic(session *mgo.Session, code string) (*Economic, error) {
	var InfoGeneral StudentInformation
	var InfoEcono Economic

	MainSession := db.Cursor(session, utility.CollectionGeneral)
	EconomicSession := db.Cursor(session, utility.CollectionEconomic)

	errP := MainSession.Find(bson.M{"codigo": code}).One(&InfoGeneral)
	if errP == nil {
		errP = EconomicSession.Find(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEcono)
	}
	return &InfoEcono, errP
}

//UpdateInformationEconomic - Actualiza la informacion economica de un estudiante
//Param session		IN   "sesion de base de datos"
//Param newInfo		IN   "nueva informacion economica"
//Param code		IN   "Codigo del estudiante a consultar"
//Param keyfiledelete	OUT   "lista de nombres de archivos que no deben de estar en servidor ni historico"
//Param err		OUT   "error si es que existe"
func UpdateInformationEconomic(session *mgo.Session, newInfo *Economic, code string) (*[]string, error) {
	var InfoGeneral StudentInformation
	var InfoEcoOld Economic
	var keyfiledelete []string
	MainSession := db.Cursor(session, utility.CollectionGeneral)
	EconomicSession := db.Cursor(session, utility.CollectionEconomic)

	MainSession.Find(bson.M{"codigo": code}).One(&InfoGeneral)
	err := EconomicSession.Find(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEcoOld)
	if err == nil {
		Rescueinf(newInfo, &InfoEcoOld, &keyfiledelete)
		err = EconomicSession.Update(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}, &InfoEcoOld)
	}
	return &keyfiledelete, err
}

//GetRequiredFiles - optiene los archivos que un estudiante debe tener para poder realizar una correcta inscripcion
//Param session		IN   "sesion de base de datos"
//Param code		IN   "Codigo del estudiante a consultar"
//Param keyrequired	OUT   "lista de nombres de archivos que SI deben de estar en servidor y historico"
//Param errP		OUT   "error si es que existe"
func GetRequiredFiles(session *mgo.Session, code string) (*[]string, error) {
	var InfoGeneral StudentInformation
	var InfoEcono Economic
	MainSession := db.Cursor(session, utility.CollectionGeneral)
	EconomicSession := db.Cursor(session, utility.CollectionEconomic)

	var key = []string{"PersonasACargo", "EmpleadorOArriendo", "CondicionEspecial", "CondicionDiscapacidad", "PatologiaAlimenticia"}
	var keyrequired = []string{"FormatoInscripcion", "CartaADirectora", "CertificadoEstrato", "FotocopiaReciboServicio", "CertificadoIngresos", "ReciboUniversidad"}
	MainSession.Find(bson.M{"codigo": code}).One(&InfoGeneral)
	errP := EconomicSession.Find(bson.M{"id": InfoGeneral.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEcono)
	if errP == nil {
		if strings.Compare(InfoEcono.PersACargo, "si") == 0 {
			keyrequired = append(keyrequired, key[0])
		}
		if strings.Compare(InfoEcono.EmpleadArriendo, "si") == 0 {
			keyrequired = append(keyrequired, key[1])
		}
		if strings.Compare(InfoEcono.PobEspecial, "N") != 0 {
			keyrequired = append(keyrequired, key[2])
		}
		if strings.Compare(InfoEcono.Discapacidad, "si") == 0 {
			keyrequired = append(keyrequired, key[3])
		}
		if strings.Compare(InfoEcono.PatAlimenticia, "si") == 0 {
			keyrequired = append(keyrequired, key[4])
		}
	}
	return &keyrequired, errP
}

//UpdateStateVerificator - Actualiza el estado de un estudiante luego de ser verificado
//Param session		IN   "sesion de base de datos"
//Param code		IN   "Codigo del estudiante a consultar"
//Param info	IN    "informacion de estado, vericicador que califico y mensaje que se va a guardar"
//Param err		OUT   "error si es que existe"
func UpdateStateVerificator(session *mgo.Session, cod string, info *Economic) error {
	var InfoGeneralU StudentInformation
	var InfoEcoOldU Economic
	MainSession := db.Cursor(session, utility.CollectionGeneral)
	EconomicSession := db.Cursor(session, utility.CollectionEconomic)

	MainSession.Find(bson.M{"codigo": cod}).One(&InfoGeneralU)
	err := EconomicSession.Find(bson.M{"id": InfoGeneralU.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}).One(&InfoEcoOldU)
	VerificatorUpdate(info, &InfoEcoOldU)
	err = EconomicSession.Update(bson.M{"id": InfoGeneralU.ID, "periodo": time.Now().UTC().Year(), "semestre": utility.Semester()}, &InfoEcoOldU)

	return err
}

/* Functions Bonus*/

//Rescueinf - Actualiza la informacion del modelo y verifica archivos que deben estar
//Param newI		IN   "nueva informacion socio economica"
//Param old		IN   "antigua informacion socioeconomica"
//Param FileExists	IN    "archivos que deben estar en el historico y en el servidor"
func Rescueinf(newI *Economic, old *Economic, FileExists *[]string) {
	old.EstadoProg = newI.EstadoProg
	//replace old information
	if strings.Compare(newI.Estrato, "") != 0 {
		old.Estrato = newI.Estrato
	}
	if newI.Ingresos != 0 {
		old.Ingresos = newI.Ingresos
	}
	if strings.Compare(newI.SostePropia, "") != 0 {
		old.SostePropia = newI.SostePropia
	}
	if strings.Compare(newI.Tipoapoyo, "") != 0 {
		old.Tipoapoyo = newI.Tipoapoyo
	}
	if strings.Compare(newI.SosteHogar, "") != 0 {
		old.SosteHogar = newI.SosteHogar
	}
	if strings.Compare(newI.Nucleofam, "") != 0 {
		old.Nucleofam = newI.Nucleofam
	}
	if strings.Compare(newI.PersACargo, "") != 0 {
		if strings.Compare(newI.PersACargo, "no") == 0 {
			*FileExists = append(*FileExists, "PersonasACargo")
		}
		old.PersACargo = newI.PersACargo
	}
	if strings.Compare(newI.EmpleadArriendo, "") != 0 {
		if strings.Compare(newI.EmpleadArriendo, "no") == 0 {
			*FileExists = append(*FileExists, "EmpleadorOArriendo")
		}
		old.EmpleadArriendo = newI.EmpleadArriendo
	}
	if strings.Compare(newI.ProvBogota, "") != 0 {
		old.ProvBogota = newI.ProvBogota
	}
	if strings.Compare(newI.Ciudad, "") != 0 {
		old.Ciudad = newI.Ciudad
	}
	if strings.Compare(newI.PobEspecial, "") != 0 {
		if strings.Compare(newI.PobEspecial, "N") == 0 {
			*FileExists = append(*FileExists, "CondicionEspecial")
		}
		old.PobEspecial = newI.PobEspecial
	}
	if strings.Compare(newI.Discapacidad, "") != 0 {
		if strings.Compare(newI.Discapacidad, "no") == 0 {
			*FileExists = append(*FileExists, "CondicionDiscapacidad")
		}
		old.Discapacidad = newI.Discapacidad
	}
	if strings.Compare(newI.PatAlimenticia, "") != 0 {
		if strings.Compare(newI.PatAlimenticia, "no") == 0 {
			*FileExists = append(*FileExists, "PatologiaAlimenticia")
		}
		old.PatAlimenticia = newI.PatAlimenticia
	}
	if strings.Compare(newI.SerPiloPaga, "") != 0 {
		old.SerPiloPaga = newI.SerPiloPaga
	}
	if strings.Compare(newI.Sisben, "") != 0 {
		old.Sisben = newI.Sisben
	}
	if strings.Compare(newI.Telefono, "") != 0 {
		old.Telefono = newI.Telefono
	}
	if strings.Compare(newI.Correo, "") != 0 {
		old.Correo = newI.Correo
	}
	if strings.Compare(newI.Antiguedad, "") != 0 {
		old.Antiguedad = newI.Antiguedad
	}

	if strings.Compare(newI.TipoSubsidio, "") != 0 {
		old.TipoSubsidio = newI.TipoSubsidio
	}
}

//VerificatorUpdate - Actualiza informacion socioeconomica antigua
//Param newI	IN   "nueva informacion socio economica"
//Param old		IN   "antigua informacion socioeconomica"
func VerificatorUpdate(newI *Economic, old *Economic) {

	old.EstadoProg = newI.EstadoProg
	if strings.Compare(newI.Verificadopor, "") != 0 {
		old.Verificadopor = newI.Verificadopor
	}
	if strings.Compare(newI.Mensaje, "") != 0 {
		old.Mensaje = newI.Mensaje
	}
}
