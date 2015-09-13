#!/usr/bin/python

import json
from oauth2client import client, crypt
import sys

CLIENT_ID = '750357624031-sfnbkqo5q3bta27dr9fmsmhd3l2otj1a.apps.googleusercontent.com'

token = sys.argv[1]

try:
    idinfo = client.verify_id_token(token, CLIENT_ID)
    # If multiple clients access the backend server:
    if idinfo['aud'] != CLIENT_ID:
        raise crypt.AppIdentityError("Unrecognized client.")
    if idinfo['iss'] not in ['accounts.google.com', 'https://accounts.google.com']:
        raise crypt.AppIdentityError("Wrong issuer.")
except crypt.AppIdentityError as e:
    # Invalid token
    sys.stderr.write('Invalid Token: ' + str(e) + '\n')
    sys.exit(1)

print json.dumps(idinfo)