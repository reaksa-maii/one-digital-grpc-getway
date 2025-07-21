package onegrpcgetway
import "log"

func LogError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}
func LogInfo(message string) {
	log.Printf("Info: %s", message)
}
func LogDebug(message string) {
	log.Printf("Debug: %s", message)
}
func LogWarning(message string) {
	log.Printf("Warning: %s", message)
}
func LogFatal(err error) {
	if err != nil {
		log.Fatalf("Fatal error: %v", err)
	}
}