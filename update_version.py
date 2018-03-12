#!/usr/bin/env python3

# Project values
project_go = "dataset.go"
codemeta_json = "codemeta.json"

#
# No changes below this line
#
import sys
import os
import json

def inc_patch_no(v = "0.0.0"):
    """inc_patch_no takes a symvar and increments the right most value in the dot touple"""
    parts = v.split(".")
    if len(parts) == 3:
        #major_no = parts[0]
        #minor_no = parts[1]
        patch_no = int(parts[2])
        patch_no += 1
        parts[2] = str(patch_no)
        return ".".join(parts)

    else:
        return v

def update_codemeta_json(codemeta_json, current_version, next_version):
    with open(codemeta_json, mode = "r", encoding = "utf-8") as f:
        src = f.read()
    meta = json.loads(src)
    meta["version"] = next_version
    downloadURL = meta["downloadUrl"]
    meta["downloadUrl"] = downloadURL.replace(current_version, next_version)
    src = json.dumps(meta)
    print(f"updating {codemeta_json} version from {current_version} to {next_version}")

    with open(codemeta_json, mode = "w", encoding = "utf-8") as f:
        f.write(src)
    return True

def update_project_go(project_go, current_version, next_version):
    current_version = f"v{current_version}"
    next_version = f"v{next_version}"

    print(f"updating {project_go} Version from {current_version} to {next_version}")
    print("WARNING: update_project_go not implemented")
    return True

def usage(app_name):
    app_name = os.path.basename(app_name)
    print(f"""
USAGE: {app_name} OPTIONS

SYNOPSIS

{app_name} shows or sets the proposed new value for a version number.
By defaut it proposes a increment in the patch no of a symvar string.
If the -y, --yes option is included it will commit the change in patch
number to the codemeta.json and project's go file.

OPTIONS

    --set VALUE      explicitly set the value of the new version string
    -y, --yes        commit the changes proposed to the Codemeta and Go file.
""")

#
# Main processing
#
def main(args):
    if ("-h" in args) or ("-help" in args) or ("--help" in args):
        usage(args[0])
        sys.exit(0)
    current_version = ""
    next_version = ""
    meta = {}
    with open(codemeta_json,"r") as f:
        src = f.read()
        meta = json.loads(src)

    current_version = meta["version"]

    if ("--set" in args):
        i = args.index("--set")
        i += 1
        if len(args) < i:
            print("Missing new version number after set", args)
            sys.exit(1)
        next_version = args[i]
        if next_version[0] == "v":
            next_version = next_version[1:]
    else:
        next_version = inc_patch_no(current_version)

    if ("--yes" in args) or ("-y" in args):
        ok = update_codemeta_json(codemeta_json, current_version, next_version)
        if ok == False:
            sys.exit(1)
        ok = update_project_go(project_go, current_version, next_version)
        if ok == False:
            sys.exit(1)
        sys.exit(0)
    else:
        print("current version:", current_version)
        print("proposed version:", next_version)

if __name__ == "__main__":
    main(sys.argv[:])
