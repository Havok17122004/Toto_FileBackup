# Toto_FileBackup
I wrote this code for a PClub project of creating a CLI-based backup tool in Go.

Welcome to toto. toto is a CLI backup tool produced using golang and Cobra.


This program is a big failure. The person who wrote it's code is dumb, but he surely 
did a lot of hard work to learn development and implement it in his first-ever project
This is an incomplete project but I hope this will have future updates in the future!

Find the documentation on https://github.com/Havok17122004/Toto_FileBackup.git

![toto-high-resolution-logo](https://github.com/Havok17122004/Toto_FileBackup/assets/148974367/03afafee-9196-46a7-9243-710e9f162dc0)


## Commands ->

#### 1) copyall -
  This copies the entire contents of a folder including subfolders and subfiles directly into the destination folder named 'backup'. It may replace any existing files present in the 'backup' directory, if present.
	If the source directory is: "D:\test"
	Destination directory is: "D:\My Tickets"
	The command can be run as: "D:\test" "D:\My Tickets"

 #### 2) copy-
   This copies only a specified file from the source directory to a destination directory
   If the source directory is: "D:\test"
   Destination directory is: "D:\My Tickets"
   The command can be run as : .\toto copy "D:\test" "D:\My Tickets"

   
Open Bug: Cannot correctly identify which files are to be modified and which are not.
