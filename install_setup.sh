#clear bin
rm -rf ./bin/*

#copy resource files
mkdir bin/dev
touch bin/dev/.gitignore
mkdir bin/downloads
mkdir bin/downloads/custom
touch bin/downloads/custom/.gitignore
mkdir bin/downloads/design_recommend_files
touch bin/downloads/design_recommend_files/.gitignore
mkdir bin/downloads/custom_processed
touch bin/downloads/custom_processed/.gitignore
mkdir bin/downloads/report
touch bin/downloads/report/.gitignore
mkdir bin/downloads/heads
touch bin/downloads/heads/.gitignore
mkdir bin/snapshot
touch bin/snapshot/.gitignore
mkdir bin/reports
touch bin/reports/.gitignore
mkdir bin/test
touch bin/test/.gitignore
mkdir bin/file_server_files
touch bin/file_server_files/.gitignore
mkdir bin/downloads/head
touch bin/downloads/head/.gitignore
mkdir bin/bin_java
touch bin/bin_java/.gitignore
cp build_pixie_server.txt bin/README.txt
cp -r ./downloads/*.png bin/downloads/
cp -r ./downloads/heads/*.png bin/downloads/heads/
cp -r ./scripts bin/
cp -r ./web_res bin/
cp -r ./cert bin/
cp -r ./gm bin/
# cp -r ./gm_pages bin/
# cp -r ./gm_pages_en bin/
# cp -r ./gv_pages bin/
# cp -r ./gv_pages_en bin/
cp -r ./web3py bin/
# cp -r ./docs bin/
rm -rf odpscmd/plugins/dship/sessions
cp -r ./odpscmd bin/
cp -r ./GFTP bin/
cp process_file.sh bin/
# cp process_file_all.sh bin/
cp update_version.sh bin/
# cp process_file_pvr.sh bin/
# cp process_file_pvr_all.sh bin/
# cp bin_java/RSAUtil.class bin/bin_java/
# cp bin_java/RSA.class bin/bin_java/
# cp bin_java/OppoCallbackSignUtil.class bin/bin_java/
# cp bin_java/MhrRSA.class bin/bin_java/

rm -rf bin/GFTP/clothes/custom_*
rm -rf bin/GFTP/icons/custom_*
rm -rf bin/GFTP/output/custom_*
rm -rf bin/GFTP/output_pvr/custom_*
rm -rf bin/GFTP/clothes/tmp_*
rm -rf bin/GFTP/icons/tmp_*
rm -rf bin/GFTP/output/tmp_*
rm -rf bin/GFTP/output_pvr/tmp_*
rm -rf bin/GFTP/*.deb
rm -rf bin/GFTP/.DS_Store
rm -rf bin/web_res/.DS_Store
