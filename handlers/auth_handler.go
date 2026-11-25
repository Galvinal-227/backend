package handlers

import (
"net/http"

"github.com/gin-gonic/gin"

"pkl-attendance/backend/services"
"pkl-attendance/backend/utils"
)

type AuthHandler struct {
svc *services.AuthService
}

func NewAuthHandler(svc *services.AuthService) *AuthHandler {
return &AuthHandler{svc: svc}
}

type loginRequest struct {
Username string `json:"username" binding:"required,min=3,max=64"`
Password string `json:"password" binding:"required,min=6,max=128"`
}

func (h *AuthHandler) Login(c *gin.Context) {
var req loginRequest
if err := c.ShouldBindJSON(&req); err != nil {
utils.Error(c, http.StatusBadRequest, "invalid input: "+err.Error())
return
}
result, err := h.svc.Login(req.Username, req.Password)
if err != nil {
utils.Error(c, http.StatusUnauthorized, err.Error())
return
}
utils.Success(c, http.StatusOK, result)
}

func (h *AuthHandler) Logout(c *gin.Context) {
jtiVal, exists := c.Get("jti")
if !exists {
utils.Error(c, http.StatusUnauthorized, "no session")
return
}
expVal, _ := c.Get("exp")
jti, _ := jtiVal.(string)
exp, _ := expVal.(int64)
if err := h.svc.Logout(jti, exp); err != nil {
utils.Error(c, http.StatusInternalServerError, "failed to logout")
return
}
utils.Success(c, http.StatusOK, gin.H{"message": "logged out"})
}

func (h *AuthHandler) Me(c *gin.Context) {
idVal, _ := c.Get("user_id")
id, _ := idVal.(uint)
user, err := h.svc.Me(id)
if err != nil {
utils.Error(c, http.StatusNotFound, "user not found")
return
}
utils.Success(c, http.StatusOK, user)
}

type changePasswordRequest struct {
OldPassword string `json:"old_password" binding:"required,min=6"`
NewPassword string `json:"new_password" binding:"required,min=6,max=128"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
var req changePasswordRequest
if err := c.ShouldBindJSON(&req); err != nil {
utils.Error(c, http.StatusBadRequest, "input tidak valid: "+err.Error())
return
}
if err := h.svc.ChangePassword(c.GetUint("user_id"), req.OldPassword, req.NewPassword); err != nil {
utils.Error(c, http.StatusBadRequest, err.Error())
return
}
utils.Success(c, http.StatusOK, gin.H{"message": "password berhasil diubah"})
}