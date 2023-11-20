**How to Use**

1. clone the repository

`$ git clone github.com/mohammed-maher/WistiaCourseDownloader`

2. cd to the directory & build the binary

`$ cd WistiaCourseDownloader && go build main.go`

3. create downloads folder

`$ mkdir downloads`

4. construct your input.csv file

the input.csv file must be populated before executing the program

it has the following format
Section,Name,id

Section: the name of the module, note that you MUST create a directory inside downloads directory for every unique Section name.
Name: The title of the video to be saved, you don't need the mp4 extension as it will be added to each name
id: the id of the video on Wistia, the ID is usually a 10 chars string you
to find it, right click on the video and click on Copy link and thumbnail, then paste it anywhere and find "video=xxxxxxxxxx" add only the xxxxxxxxxx as the video id


5. execute the program

`$ ./main`

wait until the download completes, you should be able to find the videos in the downloads folder