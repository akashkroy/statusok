package notify

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

//Diffrent types of clients to deliver notifications
type NotificationTypes struct {
	MailNotify MailNotify      `json:"mail"`
	Mailgun    MailgunNotify   `json:"mailGun"`
	Slack      SlackNotify     `json:"slack"`
	Http       HttpNotify      `json:"httpEndPoint"`
	Dingding   DingdingNotify  `json:"dingding"`
	Pagerduty  PagerdutyNotify `json:"pagerduty"`
}

type ResponseTimeNotification struct {
	Url                  string
	RequestType          string
	ExpectedResponsetime int64
	MeanResponseTime     int64
}

type ErrorNotification struct {
	Url          string
	RequestType  string
	ResponseBody string
	Error        string
	OtherInfo    string
}

var (
	errorCount        = 0
	notificationsList []Notify
)

type Notify interface {
	GetClientName() string
	Initialize() error
	SendResponseTimeNotification(notification ResponseTimeNotification) error
	SendErrorNotification(notification ErrorNotification) error
}

//Add notification clients given by user in config file to notificationsList
func AddNew(notificationTypes NotificationTypes) {

	v := reflect.ValueOf(notificationTypes)

	for i := 0; i < v.NumField(); i++ {
		notifyString := fmt.Sprint(v.Field(i).Interface().(Notify))
		//Check whether notify object is empty . if its not empty add to the list
		if !isEmptyObject(notifyString) {
			notificationsList = append(notificationsList, v.Field(i).Interface().(Notify))
		}
	}

	if len(notificationsList) == 0 {
		println("No clients Registered for Notifications")
	} else {
		println("Initializing Notification Clients....")
	}

	for _, value := range notificationsList {
		initErr := value.Initialize()

		if initErr != nil {
			println("Notifications : Failed to Initialize ", value.GetClientName(), ".Please check the details in config file ")
			println("Error Details :", initErr.Error())
		} else {
			println("Notifications :", value.GetClientName(), " Intialized")
		}

	}
}

//Send response time notification to all clients registered
func SendResponseTimeNotification(responseTimeNotification ResponseTimeNotification) {

	for _, value := range notificationsList {
		err := value.SendResponseTimeNotification(responseTimeNotification)

		//TODO: exponential retry if fails ? what to do when error occurs ?
		if err != nil {

		}
	}
}

//Send Error notification to all clients registered
func SendErrorNotification(errorNotification ErrorNotification) {

	for _, value := range notificationsList {
		err := value.SendErrorNotification(errorNotification)

		//TODO: exponential retry if fails ? what to do when error occurs ?
		if err != nil {

		}
	}
}

//Send Test notification to all registered clients .To make sure everything is working
func SendTestNotification() {

	println(":white_check_mark: Statusok Test Notification")
}

func validateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func isEmptyObject(objectString string) bool {
	objectString = strings.Replace(objectString, "0", "", -1)
	objectString = strings.Replace(objectString, "map", "", -1)
	objectString = strings.Replace(objectString, "[]", "", -1)
	objectString = strings.Replace(objectString, " ", "", -1)
	objectString = strings.Replace(objectString, "{", "", -1)
	objectString = strings.Replace(objectString, "}", "", -1)

	if len(objectString) > 0 {
		return false
	} else {
		return true
	}
}

//A readable message string from responseTimeNotification
func getMessageFromResponseTimeNotification(responseTimeNotification ResponseTimeNotification) string {

	message := fmt.Sprintf(":warning: Increased API response time on %v %v %v ms / %v ms",
		responseTimeNotification.Url, responseTimeNotification.RequestType, responseTimeNotification.MeanResponseTime, responseTimeNotification.ExpectedResponsetime)

	return message
}

//A readable message string from errorNotification
func getMessageFromErrorNotification(errorNotification ErrorNotification) string {

	message := fmt.Sprintf(":red_circle: Error sending requests to URL: %v %v\n"+
		"Error Message: %v Response: %v\nOther Info:%v",
		errorNotification.Url, errorNotification.RequestType, errorNotification.Error, errorNotification.ResponseBody, errorNotification.OtherInfo)

	return message
}
