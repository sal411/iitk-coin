Controllers and Models packages are for structuring the prorgam,
which could not be done

Errors : File hee nhi import ho rhi thi bc, kya randi cheez hai go, 
local package import krne mein itna chod bc, dimaag ka dahi ho gya,
jitna code likhne mein time laga usse 4 guna sirf ye figure out krne mein ki ye import kyu nhi ho tha
Sorry for that rant, Lets move on 

I will look into this later, so for now, I have written the whole program in single file

But the flow is as follows 

user
    user.go      -------- this has fuctions for user struct, connectDB, add user, newuser
    userdata.go  -------- has user struct
main

main has several hadlerFunc
/signup
/login
/secretpage
/

/signup 
    it runs createUser which uses userData struct to create user in DB

/login 
    this first finds the hashedpassword saved in db for the corresponding rollno, and return error accordingly,
            then creates a JWT token and sets it in cookie, valid for 10 minutes 

/secretpage 
    this first takes the token  from cookie and verfies it, and if matched, presents the information, else rejects it 

