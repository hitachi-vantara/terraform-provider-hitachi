package terraform

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	commonlog "terraform-provider-hitachi/hitachi/common/log"
	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var syncChangeUserPasswordOperation = &sync.Mutex{}

// Resource for changing user password
func ResourceVssbChangeUserPassword() *schema.Resource {
	return &schema.Resource{
		Description:   "VOS Block: Change Storage User Password.",
		CreateContext: resourceVssbChangeUserPasswordCreate,
		UpdateContext: resourceVssbChangeUserPasswordUpdate,
		DeleteContext: resourceVssbChangeUserPasswordDelete,
		ReadContext:   resourceVssbChangeUserPasswordRead,
		Schema:        schemaimpl.ResourceVssbChangeUserPasswordSchema,
		CustomizeDiff: validatePasswordChangeInputs,
	}
}

func resourceVssbChangeUserPasswordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()
	syncChangeUserPasswordOperation.Lock()
	defer syncChangeUserPasswordOperation.Unlock()

	log.WriteInfo("starting change user password")
	userInfo, err := impl.ChangeVssbUserPassword(d)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	userInfoMap := impl.ConvertVssbStorageUserToSchema(userInfo)
	if err := d.Set("user_info", []interface{}{userInfoMap}); err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	d.SetId(userInfo.UserId)
	d.Set("status", "Password changed successfully")
	log.WriteInfo("password changed successfully")
	return nil
}

func resourceVssbChangeUserPasswordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceVssbChangeUserPasswordCreate(ctx, d, m)
}

func resourceVssbChangeUserPasswordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceVssbChangeUserPasswordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

func validatePasswordChangeInputs(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return validatePasswordChangeInputsLogic(d)
}

type minimalDiff interface {
	Get(string) interface{}
	GetOk(string) (interface{}, bool)
}

// testable
func validatePasswordChangeInputsLogic(diff minimalDiff) error {
	currentPassword := diff.Get("current_password").(string)
	newPassword := diff.Get("new_password").(string)
	userID := diff.Get("user_id").(string)

	userIDRegex := regexp.MustCompile(`^[-A-Za-z0-9!#$%&'.@^_{}~]{5,255}$`)
	passwordRegex := regexp.MustCompile(`^[-A-Za-z0-9!#$%&"'()*+,./:;<>=?@[\]\\^_` + "`" + `{|}~]{1,256}$`)

	if userID == "" {
		return fmt.Errorf("missing user_id")
	}
	if len(userID) < 5 || len(userID) > 255 || !userIDRegex.MatchString(userID) {
		return fmt.Errorf("user_id must be 5 to 255 valid characters")
	}

	if currentPassword == newPassword {
		return fmt.Errorf("new_password must be different from current_password")
	}

	if len(currentPassword) < 1 || len(currentPassword) > 256 || !passwordRegex.MatchString(currentPassword) {
		return fmt.Errorf("current_password must be 1 to 256 valid characters")
	}

	if len(newPassword) < 1 || len(newPassword) > 256 || !passwordRegex.MatchString(newPassword) {
		return fmt.Errorf("new_password must be 1 to 256 valid characters")
	}

	return nil
}
