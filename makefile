# use it in repository dir
git.log:
	@cd $(path)
	@#echo $(branch): > log-$(branch).txt 
	git ls-remote
	@#git log origin/$(branch) --pretty=format:"%h%x09%s %an%x09%ad%x09%n" --grep="PS" --since="3 weeks ago" >> src/log-$(branch).txt