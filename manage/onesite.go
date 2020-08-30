package manage

import "jinliconfig/class"

func OneCreateSite(BASEPATH string, newDomain string, category string) {
	class.ExecLinuxCommand("wget -O  /tmp/latest.zip https://release.jinli.plus/app/" + category + "/latest.zip")
	class.ExecLinuxCommand("cd " + BASEPATH + "code/" + newDomain + "/" + " && rm -rf *")
	class.ExecLinuxCommand("mv /tmp/latest.zip " + BASEPATH + "code/" + newDomain + "/")
	class.ExecLinuxCommand("cd " + BASEPATH + "code/" + newDomain + "/" + " && unzip latest.zip")
	class.ExecLinuxCommand("rm -f " + BASEPATH + "code/" + newDomain + "/" + "latest.zip")
	class.ExecLinuxCommand("cd " + BASEPATH + "code/" + newDomain + " && chown -R 10000:10000 *")
}
