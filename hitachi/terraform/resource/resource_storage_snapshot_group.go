package terraform

import (
    "context"
    "fmt"
    "sync"

    commonlog "terraform-provider-hitachi/hitachi/common/log"
    impl "terraform-provider-hitachi/hitachi/terraform/impl"
    schemaimpl "terraform-provider-hitachi/hitachi/terraform/schema"

    "github.com/hashicorp/terraform-plugin-sdk/v2/diag"
    "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Reusing a mutex to ensure serial execution of storage operations if required
var syncSnapshotGroupOperation = &sync.Mutex{}

func ResourceVspSnapshotGroup() *schema.Resource {
    return &schema.Resource{
        Description:   `Vsp Snapshot Group Resource: Manages snapshots in units of snapshot groups (Consistency Groups) on Hitachi VSP storage.`,
        CreateContext: resourceVspSnapshotGroupCreate,
        ReadContext:   resourceVspSnapshotGroupRead,
        UpdateContext: resourceVspSnapshotGroupUpdate,
        DeleteContext: resourceVspSnapshotGroupDelete,
        Schema:        schemaimpl.ResourceVspSnapshotGroupSchema(),
        CustomizeDiff: resourceVspSnapshotGroupCustomizeDiff,
    }
}

func resourceVspSnapshotGroupCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    syncSnapshotGroupOperation.Lock()
    defer syncSnapshotGroupOperation.Unlock()

    // Typically calls the logic to perform the initial group action (e.g., create/read)
    return impl.ResourceVspSnapshotGroupApply(d)
}

func resourceVspSnapshotGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    return impl.ResourceVspSnapshotGroupRead(d)
}

func resourceVspSnapshotGroupUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    syncSnapshotGroupOperation.Lock()
    defer syncSnapshotGroupOperation.Unlock()

    // Handles state transitions like split, resync, or restore
    return impl.ResourceVspSnapshotGroupApply(d)
}

func resourceVspSnapshotGroupDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
    syncSnapshotGroupOperation.Lock()
    defer syncSnapshotGroupOperation.Unlock()

    // Note: Deleting a group resource usually involves deleting the pairs within or the group definition
    return impl.ResourceVspSnapshotGroupDelete(d)
}

func resourceVspSnapshotGroupCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
    log := commonlog.GetLogger()

    // --- 1. IMMUTABILITY CHECKS ---
    // If the resource already exists, prevent changing core identifiers
    if d.Id() != "" {
        immutableFields := []string{
            "serial",
            "snapshot_group_name",
        }

        for _, field := range immutableFields {
            if d.HasChange(field) {
                old, new := d.GetChange(field)

                // Allow if transitioning to a zero value (destruction phase)
                if isZeroValue(new) {
                    continue
                }

                err := fmt.Errorf("%s is immutable: cannot change from %v to %v. Recreate resource for a different group", field, old, new)
                log.WriteError("TFError| SnapshotGroup Immutability Violation: %v", err)
                return err
            }
        }
    }

    // Mark output as computed so it refreshes after the action
    d.SetNewComputed("snapshot_group")

    return nil
}
