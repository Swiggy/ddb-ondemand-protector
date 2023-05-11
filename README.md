# DoP - DynamoDB OnDemand Protector
DynamoDB OnDemand Protector is a pipeline that helps to put cost guardrails on a system if it uses on-demand DynamoDB tables.
## Motivation
While on-demand DynamoDB tables allow the pay per usage concept, the RCU/WCU limits are on an account level and not at a table level. This might result in a huge costs on the system if there is a very high RCU/WCU on a table for a prolonged duration. In our case, this happened as part of data load job to pre-populate a table. Since the limits imposed on the account were very high, the on-demand table scaled out to the maximum capacity incurring huge costs. If you have a similar case [Imports from s3 to ddb](https://aws.amazon.com/blogs/database/amazon-dynamodb-can-now-import-amazon-s3-data-into-a-new-table/) might help. With this pipeline, we have created guardrails at a table level to limit the RCU/WCU whenever the cost threshold reaches for a table.

For example, let's assume the RCU and WCU limit to be 40000 on the account level for on-demand mode. This might have been decided to support one of the critical tables.
Lets find the cost incurred on the writes to the table for 1 hour: <br />
Assumed size of item - 1KB <br />
Region - ap-southeast-1 <br />
1 WCU is required for each request, as for items up to 1 KB in size, one WCU can perform one standard write request per second [reference](https://aws.amazon.com/dynamodb/pricing/provisioned/) <br />
Cost incurred for 1 hour writes -  40000 * 60 * 60 * $0.0000014231 = $204.9264 per hour <br />
Now, this may be justified for the critical table but let's say there is another table that is on-demand and is used for non-critical cases wherein you want to limit the cost consumption to ~$102 per hour. In this case, if the non-critical table consumes all the WCUs the cost will shoot up to $204 per hour and the other tables have to compete for WCUs as the limits are at account level. <br />
To solve this issue, we can put limits on table level using the dynamodb ondemand protector. The first step in this is to identify the threshold consumption you want to put for your table. In this case, if you want to limit your table cost limit to ~$102 per hour for writes, your WCU threshold should be 20000. You can find your WCU RCU threshold using this [calculator](https://calculator.aws/#/addService/DynamoDB). You can configure the RCU WCU to get approximate cost and accordingly tune your thresholds.

Once the thresholds are defined, identify the RCU WCU you would want to put once the thresholds are reached. Note that these thresholds will be for a provisioned capacity. That is, let's say you have defined RCU WCU both to be 1 once the threshold is breached. This would mean that the table will now be converted to provisioned mode with the given RCU WCU.
Once you have identified the numbers, configure them in the ondemand protector and deploy the same for your table.

## Architecture
The pipeline is based on AWS step functions that will switch the mode of an on-demand table to provisioned when it exceeds a predefined RCU/WCU threshold. 
The mode shall be switched back to on-demand after a specified time. AWS allows switching between the modes once every 24 hours, the default time for the switch is taken to be 24 hours.  
![Alt text](assets/architecture.png?raw=true "Title")  
This [blog](https://bytes.swiggy.com/how-to-limit-autoscaling-in-on-demand-dynamodb-tables-c57e20cbbbcf) can be referred for more details.
## How to Use?
The mode of deployment of this pipeline is via [AWS SAM](https://aws.amazon.com/serverless/sam/).  
Create the sam template to deploy the resources using following command: 
```python create_template.py [-h] [-l LAMBDA_EXECUTE_ROLE] [-n SNS_NOTIFY] [-s STEP_EXECUTE_ROLE] [-a SNS_ALERT] [-d TABLE_NAME] [-r THRESHOLD] [-t STEP_TRIGGER_ROLE] [-o OUTPUT_FILE]```
where:
- LAMBDA_EXECUTE_ROLE: IAM role to execute lambda which will switch the ddb table
- SNS_NOTIFY: SNS topic which will send notifications for provisionibg mode switch
- STEP_EXECUTE_ROLE: IAM role to execute step function
- SNS_ALERT: SNS topic to send alerts when threshold is breached and switch is triggered
- TABLE_NAME: DynamoDB table name on which the limiting is to be applied
- THRESHOLD: Threshold for RCU
- STEP_TRIGGER_ROLE: IAM role to trigger step function 
- OUTPUT_FILE: output yaml sam template

On running the above script, sam template at location given in the above command is ready to use.

The Cloudformation template has the following resources:
- DdbSwitch - The lambda function responsible for the switch
- DdbStep - The step function orchestrating the entire flow
- Ddbalarm - Alarm configured on the RCU/WCU of the table on the breach of which
  you want to limit the RCU of the table.
- DdbSwitchRule - The rule to trigger the step function. This rule will be activated once the defined thresholds are reached. The automatic action taken by the limiter is to switch the mode of the on-demand table to provisioned with the defined RCU WCU for 24 hours. After 24 hours, the table will be switched back o on-demand.

The current template is configured for RCU. To configure WCU, replicate the Ddbalarm and DdbSwitchRule for WriteCapacityUnits.

### Steps to deploy
1. Install SAM CLI
2. Run ```sam build -t <output file>```
3. Run ```sam deploy```
