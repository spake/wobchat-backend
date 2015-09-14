#!/usr/bin/python

import json
from oauth2client import client, crypt
import sys

CLIENT_ID = '783279836221-m71iri9830ptguifn0apfbsnj22pfeel.apps.googleusercontent.com'

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
    print 'Invalid Token: ' + str(e)
    sys.exit(1)

print json.dumps(idinfo)