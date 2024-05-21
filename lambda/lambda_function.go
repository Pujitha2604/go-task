package main
 
import (
    "context"
    "encoding/json"
    "fmt"
    "log"
 
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
)
 
// ECRDetail defines the structure for the image tag details in the ECR event.
type ECRDetail struct {
    ImageTags []string `json:"imageTags"`
}
 
// ECREvent defines the structure of the ECR event payload.
type ECREvent struct {
    Detail ECRDetail `json:"detail"`
}
 
// Handler function to process the CloudWatch Event.
func handler(ctx context.Context, event events.CloudWatchEvent) error {
    // Unmarshal the event detail to extract the ECR event information.
    var ecrEvent ECREvent
    if err := json.Unmarshal(event.Detail, &ecrEvent); err != nil {
        return fmt.Errorf("error unmarshalling event detail: %v", err)
    }
 
    // Check if image tags are present in the event detail.
    if len(ecrEvent.Detail.ImageTags) == 0 {
        return fmt.Errorf("no image tags found in event detail")
    }
 
    // Extract the first image tag from the list.
    imageTag := ecrEvent.Detail.ImageTags[0]
    log.Printf("Received image tag: %s", imageTag)
 
    // Further processing with the image tag can be done here.
 
    return nil
}
 
func main() {
    // Start the Lambda function.
    lambda.Start(handler)
}
 