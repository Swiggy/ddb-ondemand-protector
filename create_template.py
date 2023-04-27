#!/usr/bin/env python3
import argparse

def create_template(input, exec_role_lambda, topic_name, exec_role_step, alert_topic_name, table, threshold,
                    role_trigger_step, output):
    with open(input, 'r') as file:
        data = file.read()
        data = data.replace("<IAM execution role for lambda>", exec_role_lambda)
        data = data.replace("<SNS topic to send switch notifications>", topic_name)
        data = data.replace("<IAM role for executing state machine>", exec_role_step)
        data = data.replace("<SNS topic to send alerts>", alert_topic_name)
        data = data.replace("<Table>", table)
        data = data.replace("<Threshold for RCU>", "'" + threshold + "'")
        data = data.replace("<IAM role to trigger the step function>", role_trigger_step)
    with open(output, 'w') as file:
        file.write(data)


if __name__ == "__main__":
    # Initialize parser
    parser = argparse.ArgumentParser()
    parser.add_argument("-l", "--lambda_execute_role", help="IAM execution role for lambda")
    parser.add_argument("-n", "--sns_notify", help="SNS topic to send notifications")
    parser.add_argument("-s", "--step_execute_role", help="IAM execution role for executing step function")
    parser.add_argument("-a", "--sns_alert", help="SNS topic to send alerts/Pagerduty notifications")
    parser.add_argument("-d", "--table_name", help="Ddb table you want to limit the RCU")
    parser.add_argument("-r", "--threshold", help="Threshold for RCU")
    parser.add_argument("-t", "--step_trigger_role", help="IAM role to trigger step function")
    parser.add_argument("-o", "--output_file", help="Output file name")

    args = parser.parse_args()
    input = "sam/cu-limiter-ddb.yaml"
    exec_role_lambda = args.lambda_execute_role
    topic_name = args.sns_notify
    exec_role_step = args.step_execute_role
    alert_topic_name = args.sns_alert
    table = args.table_name
    threshold = args.threshold
    role_trigger_step = args.step_trigger_role
    output = args.output_file
    create_template(input, exec_role_lambda, topic_name, exec_role_step, alert_topic_name, table, threshold,
                    role_trigger_step, output)
