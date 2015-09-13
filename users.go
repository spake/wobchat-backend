package main

import "os/exec"

type User struct {
    ID      string  `gorm:"primary_key"`
    Name    string
    Email   string
    Picture string
}

func verifyIdToken(token string) string {
  verifierCmd := exec.Command("./verify_token.py", token)
  idinfo, err := verifierCmd.Output()
  if err != nil {
    // unable to verify the token
    return ""
  }
  // for now just return this as a json string until we know
  // what it looks like
  return string(idinfo)
}