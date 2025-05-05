package users

// func Login(c *gin.Context) {
// 	var authInput userSchemas.AuthInput

// 	if err := c.ShouldBindJSON(&authInput); err != nil {
// 		utils.Error(c.Writer, http.StatusBadRequest, utils.ErrInvalidInput, err.Error())
// 		return
// 	}

// 	var userFound models.User
// 	initializers.DB.Where("email=?", authInput.Email).Find(&userFound)

// 	if userFound.ID == uuid.Nil {
// 		utils.Error(c.Writer, http.StatusNotFound, "User not found", nil)
// 		return
// 	}

// 	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(authInput.Password)); err != nil {
// 		utils.Error(c.Writer, http.StatusUnauthorized, "Invalid password", err.Error())
// 		return
// 	}

// 	// Generate Access Token (valid for 6 hours)
// 	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"id":  userFound.ID,
// 		"exp": time.Now().Add(6 * time.Hour).Unix(),
// 	})

// 	signedAccessToken, err := accessToken.SignedString([]byte(os.Getenv("SECRET")))
// 	if err != nil {
// 		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to generate access token", err.Error())
// 		return
// 	}

// 	// Generate Refresh Token (valid for 7 days)
// 	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
// 		"id":  userFound.ID,
// 		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
// 	})

// 	signedRefreshToken, err := refreshToken.SignedString([]byte(os.Getenv("SECRET")))
// 	if err != nil {
// 		utils.Error(c.Writer, http.StatusInternalServerError, "Failed to generate refresh token", err.Error())
// 		return
// 	}

// 	// Return both tokens
// 	utils.Success(c.Writer, http.StatusOK, utils.AUTH_SUCCESS_LOGIN, gin.H{
// 		"access_token":  signedAccessToken,
// 		"refresh_token": signedRefreshToken,
// 	}, "")
// }
