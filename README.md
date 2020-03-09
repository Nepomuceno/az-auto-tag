# Azure auto tag

The goal of this proiject it is auto tag your resourcews with the a creation tag specifying wwho has created that resource. This is don looking back at the activitity log.

# Limitatiuons of the current implementation

The current implementation will always look at the last 90 days of activity log independent if it has done the same 5 minutes ago this could be easiily solved but was out of the scope for this first release. 

Current yuou can't change the tags to be created they will always be `Created-by` and `Created-by-id` The idea iot is that you would be able to specify those thrtough environment variables like the rest of the configurations.


