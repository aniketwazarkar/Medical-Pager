package fhir

import (
	"log"
)

// PatientAdapter handles transformation between ABDM/FHIR standards
// and the internal MongoDB representation.
type PatientAdapter struct {
	// configuration properties if required
}

// FetchExternalPatient is a stub for pulling an external FHIR resource
func (a *PatientAdapter) FetchExternalPatient(patientId string) (map[string]interface{}, error) {
	log.Println("FHIR STUB: fetching external FHIR Patient resource", patientId)
	return map[string]interface{}{
		"resourceType": "Patient",
		"id":           patientId,
		"status":       "active",
	}, nil
}
