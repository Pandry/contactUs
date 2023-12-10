This wants to be a simple script to handle data from forms.  
This issue came up multiple times during my career, and I was never able to find a good tool that was easy to configure and maintain (looking at you, formtools).  
This is also very much work in progress for now.

### sinks
In this program, a sink is a way to get data out of the form and send it somewhere

#### Webhook sink
#### Notion sink

#### Custom sink
To create a custom sink, you will need to change the source of this program.  
You will need to create a new sink file in the `./sinks` folder (maybe copy the `webhook.go` to start out), create the custom configuration for the webhook and write your sink.  
Once you finished, you will need to add your new sink in the `./configuration/configuration.go`'s `Sink` struct and process it in the `ParseBytes` (`./configuration/parseConfiguration.go` file).  
I know this is not very elegant, but I couldn't think of a better solution for now.  

### Forms
Forms represnt the various forms you might want to expose

Each form has a unique name, some input values (the form values), some sinks (where you want to flush the data to) and some additional config, like the redirect you want the client to follow or the captcha settings.

Keep in mind that images are out of scope for now
