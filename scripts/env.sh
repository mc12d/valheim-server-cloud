#!/bin/bash

export TF_VAR_yc_iam_token=$(yc iam create-token)
export AWS_PROFILE=valheim

