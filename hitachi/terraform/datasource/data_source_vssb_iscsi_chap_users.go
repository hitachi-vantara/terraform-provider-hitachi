package terraform

import (
	"context"
	// "fmt"
	"strconv"
	"time"

	commonlog "terraform-provider-hitachi/hitachi/common/log"
	utils "terraform-provider-hitachi/hitachi/common/utils"

	impl "terraform-provider-hitachi/hitachi/terraform/impl"
	schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceVssbChapUsers() *schema.Resource {
	return &schema.Resource{
		Description: ":meta:subcategory:VSS Block iSCSI Target CHAP User:Obtains the information about chap users.",
		ReadContext: DataSourceVssbChapUsersRead,
		Schema:      schemaimpl.DataVssbIscsiChapUsersSchema,
	}
}

func DataSourceVssbChapUsersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log := commonlog.GetLogger()
	log.WriteEnter()
	defer log.WriteExit()

	target_chap_user, ok := d.Get("target_chap_user").(string)
	log.WriteInfo("target_chap_user %+v", target_chap_user)

	chapUserId := ""
	chapUserName := ""
	if ok {
		if utils.IsValidUUID(target_chap_user) {
			chapUserId = target_chap_user
		} else {
			chapUserName = target_chap_user
		}
	}

	target_chap_user_name, ok := d.Get("target_chap_user_name").(string)
	log.WriteInfo("target_chap_user %+v", target_chap_user_name)
	if ok {
		chapUserName = target_chap_user_name
	}

	log.WriteDebug("chapUserId: %+v\n", chapUserId)
	log.WriteDebug("chapUserName: %+v\n", chapUserName)

	if target_chap_user == "" && target_chap_user_name == "" {
		// if target_chap_user is not specified in the data file get all chap users

		chapUsers, err := impl.GetAllVssbChapUsers(d)
		if err != nil {
			return diag.FromErr(err)
		}

		log.WriteDebug("chapUsers: %+v\n", chapUsers)

		itList := []map[string]interface{}{}
		for _, cu := range chapUsers.Data {
			eachIt := impl.ConvertVssbChapUserToSchema(&cu)
			log.WriteDebug("it: %+v\n", *eachIt)
			itList = append(itList, *eachIt)
		}

		if err := d.Set("chap_users", itList); err != nil {
			log.WriteDebug("err: %v\n", err)
			return diag.FromErr(err)
		}

		d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
		log.WriteInfo("all chap users read successfully")

		return nil

	} else {
		if chapUserId != "" {
			chapUser, err := impl.GetVssbChapUserById(d)
			if err != nil {
				return diag.FromErr(err)
			}

			log.WriteDebug("chapUser: %+v\n", chapUser)
			cu := impl.ConvertVssbChapUserToSchema(chapUser)
			itList := []map[string]interface{}{
				*cu,
			}

			if err := d.Set("chap_users", itList); err != nil {
				return diag.FromErr(err)
			}

			d.SetId(chapUser.ID)
			log.WriteInfo("chap user read successfully")

		}
		if chapUserName != "" {
			chapUser, err := impl.GetVssbChapUserByName(d)
			if err != nil {
				return diag.FromErr(err)
			}

			log.WriteDebug("chapUser: %+v\n", chapUser)
			cu := impl.ConvertVssbChapUserToSchema(chapUser)
			itList := []map[string]interface{}{
				*cu,
			}

			if err := d.Set("chap_users", itList); err != nil {
				return diag.FromErr(err)
			}

			d.SetId(chapUser.ID)
			log.WriteInfo("chap user read successfully")

		}

	}
	return nil
}
