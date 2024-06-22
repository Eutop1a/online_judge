package response

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUsernameNotExist
	CodeUseNotExist
	CodeUseIDNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken
	CodeUnauthorized

	CodeErrorVerCode
	CodeExpiredVerCode
	CodeEmailExist
	CodeInternalServerError
	CodeInvalidateEmailFormat
	CodeGeneratePicError
	CodePictureError

	CodeTestCaseFormatError
	CodeProblemTitleExist
	CodeProblemTitleNotExist
	CodeProblemIDNotExist

	CodeGetUserRankError

	CodePageNotFound

	CodeUnsupportedLanguage

	CodeErrorSecret

	CodeAdminUsernameAlreadyExist
	CodeUserAlreadyRoot
	CodeUsernameAlreadyExist
	CodeObtainVerificationCode
	CodeMethodNowAllow

	CodeNeedUsername

	CodeProblemNotFound
	CodeCategoryTypeAlreadyExist
	CodeCategoryTypeDoNotExist
	CodeCategoryIsNotEmpty
	CodeDontHaveThisCategory
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:                  "success",
	CodeInvalidParam:             "请求参数错误",
	CodeUserExist:                "用户名已存在",
	CodeUseNotExist:              "用户名不存在",
	CodeUseIDNotExist:            "用户ID不存在",
	CodeInvalidPassword:          "用户名或密码错误",
	CodeServerBusy:               "服务繁忙",
	CodeNeedLogin:                "需要登录",
	CodeInvalidToken:             "无效的token",
	CodeErrorVerCode:             "验证码错误",
	CodeExpiredVerCode:           "验证码过期",
	CodeEmailExist:               "邮箱已存在",
	CodeInternalServerError:      "服务器内部错误",
	CodeInvalidateEmailFormat:    "邮箱格式错误",
	CodeGeneratePicError:         "生成图片验证码失败",
	CodePictureError:             "图片验证码错误",
	CodeTestCaseFormatError:      "测试用例格式错误",
	CodeProblemTitleExist:        "该题目标题已存在",
	CodeProblemIDNotExist:        "题目ID不存在",
	CodeProblemTitleNotExist:     "该题目标题不存在",
	CodeGetUserRankError:         "获取用户排名失败",
	CodeUsernameNotExist:         "用户名不存在",
	CodePageNotFound:             "页面不存在",
	CodeUnauthorized:             "需要管理员权限",
	CodeUnsupportedLanguage:      "不支持的编程语言",
	CodeErrorSecret:              "密钥错误",
	CodeUserAlreadyRoot:          "该用户已是管理员",
	CodeUsernameAlreadyExist:     "该用户名已存在",
	CodeObtainVerificationCode:   "需要先获取邮箱验证码",
	CodeMethodNowAllow:           "方法不允许",
	CodeNeedUsername:             "需要用户名",
	CodeProblemNotFound:          "找不到对应的题目",
	CodeCategoryTypeAlreadyExist: "该分类名称已经存在",
	CodeCategoryTypeDoNotExist:   "该分类ID不存在",
	CodeCategoryIsNotEmpty:       "分类下题目列表非空",
	CodeDontHaveThisCategory:     "分类不存在或对应的题目列表为空",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
