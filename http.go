package main

import (
	"contactUs/configuration"
	"fmt"
	"net"
	"net/http"
	"net/url"
)

const idFieldName = "id"

func setupHttp() {
	http.HandleFunc("/submit", handleSubmission)
	http.ListenAndServe(":8888", nil)
}

func handleSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		returnError(w, 405, "Request was not a POST")
		return
	}
	if err := r.ParseForm(); err != nil {
		returnError(w, 400, "Unable to parse request's body")
		return
	}
	if r.PostForm[idFieldName] == nil || len(r.PostForm[idFieldName]) < 1 {
		returnError(w, 404, "Form ID was not found in the body")
		return
	}

	fromId := r.PostForm[idFieldName][0]
	handleForm(fromId, r.PostForm, r, w)
	return
}

func returnError(w http.ResponseWriter, errorCode int, errorStr string) {
	w.WriteHeader(errorCode)
	w.Write([]byte(errorStr))
}

func handleForm(formId string, submittedFields url.Values, r *http.Request, w http.ResponseWriter) {
	form, ok := configuration.Config.Forms[formId]
	if !ok {
		returnError(w, 404, "Could not find referenced form")
		return
	}
	fmt.Println("Handling form", formId)

	if form.Captcha.Enabled {
		keyName := form.Captcha.Captcha.GetResponseFieldName()
		captchaResponses, ok := submittedFields[keyName]
		if !ok || len(captchaResponses) < 1 {
			returnError(w, 400, "Missing captcha parameter")
			return
		}
		captchaResponse := captchaResponses[0]
		var ip string

		if configuration.Config.ClientIPHeader != "" {
			ip = r.Header.Get(configuration.Config.ClientIPHeader)
		} else {
			var err error
			if ip, _, err = net.SplitHostPort(r.RemoteAddr); err != nil {
				fmt.Println("No sink was loaded for form", formId, ". Please check the configuration")
				returnError(w, 400, "Could not fetch IP")
				return
			}
		}

		if len(ip) < 7 {
			returnError(w, 400, "Could not get IP correctly")
			return
		}

		if !form.Captcha.Captcha.VerifyCaptcha(captchaResponse, ip) {
			returnError(w, 400, "Captcha response seems incorect")
			return
		}
	}

	formValues := make(map[string]string)
	for _, field := range form.Inputs {
		if val, ok := submittedFields[field]; ok && len(val) > 0 {
			formValues[field] = val[0]
		}
	}

	if len(form.Sinks) < 1 {
		fmt.Println("No sink was loaded for form", formId, ". Please check the configuration")
	}
	failedSinkAttempts := 0
	totalSinkAttempts := 0
	for _, sinkName := range form.Sinks {
		sink, ok := configuration.Config.Sinks[sinkName]
		if !ok {
			returnError(w, 500, "Could not find referenced form in map")
			return
		}

		if len(sink.ActiveSinksList) < 1 {
			fmt.Println("No provider found in sink", sinkName, "for form", formId)
			continue
		}

		totalSinkAttempts++
		for _, subSink := range sink.ActiveSinksList {
			if err := subSink.Sink(formValues); err != nil {
				failedSinkAttempts++
				fmt.Println("Got error while processing sink ", subSink.Name(), " with inputs: ", submittedFields, ". Error: "+err.Error())
			}
		}

	}
	if totalSinkAttempts == failedSinkAttempts {
		returnError(w, 500, "There was an issue accepting your form. Please try again later")
		return
	}

	if form.Redirect != "" {
		w.Header().Add("Location", form.Redirect)
		w.WriteHeader(301)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte("Your submission was received!"))
}
