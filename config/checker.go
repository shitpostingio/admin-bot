package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/shitpostingio/admin-bot/config/structs"
)

// CheckMandatoryFields uses reflection to see if there are
// mandatory fields with zero value
func CheckMandatoryFields(isReload bool, config structs.Config) error {
	return checkStruct(isReload, reflect.TypeOf(config), reflect.ValueOf(config))
}

// checkWebhookConfig perform checks and sets default values for
// webhook configuration details
func checkWebhookConfig(config *structs.WebhookConfiguration) {

	// To use a reverse proxy, we must bind to localhost.
	if config.ReverseProxy {

		config.IP = "127.0.0.1"
		if !structs.IsStandardPort(config.ReverseProxyPort) {
			log.Fatal("checkWebhookConfig: cannot use non-standard reverse proxy port")
		}

	} else {

		if !structs.IsStandardPort(config.Port) {
			log.Fatal("checkWebhookConfig: cannot use non-standard port when ReverseProxy is disabled")
		}

	}

	if config.Domain == "" { // Domain not set
		log.Fatal("Domain not set")
	} else if strings.HasPrefix(config.Domain, "http://") || strings.HasPrefix(config.Domain, "https://") {
		log.Fatal("Domain must not contain http:// or https://")
	}

	if config.TLS {
		if config.TLSCertPath == "" {
			log.Fatal("missing TLS certificate path")
		} else if config.TLSKeyPath == "" {
			log.Fatal("missing TLS key path")
		}
	}
}

// checkStruct explores structures recursively and checks if
// struct fields have a zero value
func checkStruct(isReload bool, typeToCheck reflect.Type, valueToCheck reflect.Value) (err error) {

	for i := 0; i < typeToCheck.NumField(); i++ {

		currentField := typeToCheck.Field(i)
		currentValue := valueToCheck.Field(i)

		if currentField.Type.Kind() == reflect.Struct {
			err = checkStruct(isReload, currentField.Type, currentValue)
		} else if currentField.Type.Kind() == reflect.Slice { //TODO: capire
			err = checkSlice(isReload, currentField, currentValue)
		} else {
			err = checkField(isReload, currentField, currentValue)
		}

		if err != nil {
			return
		}
	}

	return nil
}

func checkSlice(isReload bool, typeToCheck reflect.StructField, sliceToCheck reflect.Value) error {

	//only check reloadable fields if isReload is true
	if isReload {

		reloadableTagValue := typeToCheck.Tag.Get("reloadable")
		if reloadableTagValue != "true" {
			return nil
		}

	}

	typeTagValue := typeToCheck.Tag.Get("type")
	if typeTagValue == "optional" {
		return nil
	}

	if sliceToCheck.Len() == 0 {
		return fmt.Errorf("non optional slice field %s had zero length", typeToCheck.Name)
	}

	var err error
	for i := 0; i < sliceToCheck.Len(); i++ {

		item := sliceToCheck.Index(i)
		if item.Kind() == reflect.Struct {
			err = checkStruct(isReload, reflect.TypeOf(item), reflect.ValueOf(item))
		} else {

			zeroValue := reflect.Zero(item.Type())
			if item.Interface() == zeroValue.Interface() {
				return fmt.Errorf("non optional field %s had zero value at index %d", typeToCheck.Name, i)
			}

		}

		if err != nil {
			return err
		}

	}

	return nil

}

// checkField checks if a field is optional or a webhook field
// if it isn't, it checks if the field has a zero value
func checkField(isReload bool, typeToCheck reflect.StructField, valueToCheck reflect.Value) error {

	//only check reloadable fields if isReload is true
	if isReload {

		reloadableTagValue := typeToCheck.Tag.Get("reloadable")
		if reloadableTagValue != "true" {
			return nil
		}

	}

	typeTagValue := typeToCheck.Tag.Get("type")

	if typeTagValue == "optional" || typeTagValue == "webhook" {
		return nil
	}

	zeroValue := reflect.Zero(typeToCheck.Type)

	if valueToCheck.Interface() == zeroValue.Interface() {
		return fmt.Errorf("non optional field %s had zero value", typeToCheck.Name)
	}

	return nil
}
