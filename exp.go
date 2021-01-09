package main

// 	//	strunct Method, with *, allows -> Change  specific Struct by Address, change fields
// 	//init structure
// 	var Lester = person{name: "Guru", age: 21}
// 	//update fields by Pointer address struct -> fields

// 	//перегрузка методов, через struct, и interface
// 	//Lester.updatePerson( "Soma", 32)
// 	Lester.updatePerson(32)
// 	Lester.updatePerson("Jonny")

// type person struct {
// 	name string
// 	age  int
// }

// func (p *person) updatePerson(args ...interface{}) {

// 	fmt.Println(len(args))
// 	//for _, v := range data {
// 	//s := reflect.TypeOf(data)
// 	if len(args) > 1 {

// 		for i, arg := range args {
// 			// if arg.(string) != "" {
// 			// 	p.name = arg.(string)
// 			// }
// 			switch i {
// 			case 0: // name
// 				//check inside Struct index, 1 : name, 2 age..., then cast to stirng, then change value strucnt field Name -  by Pointer
// 				name, ok := arg.(string)
// 				if !ok {
// 					log.Println("error")
// 				} else {
// 					p.name = name
// 				}
// 			case 1:
// 				age, ok := arg.(int)
// 				if !ok {
// 					log.Println("error")
// 				}
// 				p.age = age
// 			default:
// 				log.Println("Wrong parametes passed")
// 			}
// 		}
// 	} else {
// 		//comapre type, then set value - variable

// 		// 	for _, arg := range args {
// 		// 		if reflect.DeepEqual(arg.(string), arg) {
// 		// 			if arg.(string) != "" {
// 		// 				p.name = arg.(string)
// 		// 			}
// 		// 		} else if reflect.DeepEqual(arg.(int), arg) {{
// 		// 			if arg.(int) > 0 {
// 		// 				p.age = arg.(int)
// 		// 			}
// 		// 		}
// 		// 	}
// 		// }
// 	}

// 	// if data.(string){
// 	// //if  (string) == reflect.TypeOf(data) {
// 	// 	fmt.Print("ds")
// 	// 	p.name = fmt.Sprintf("%s" ,data)
// 	// 	fmt.Print(p)
// 	// }else if data == reflect.Int {
// 	// 	a, _  := fmt.Printf("%v \n" ,data)
// 	// 	print(a)
// 	// 		p.age = a
// 	// 	}
// }

// // func (p *person) updatePersons(arr ...interface{}){

// // 	for _, interface := range arr {

// // for _, v := range interface {
// // 	if  v == reflect.String {
// // 		s, _ := fmt.Printf("%v" ,v)
// // 		p.name = s
// // 	}else if v == reflect.Int {
// // 			p.age = fmt.Printf("%v \n" ,v)
// // 		}

// // }
// // 	}
// // }
// // func (p *person) updatePerson(name string){
// // 	p.name = name

// // }

// if isCookie {
// 	//write cookie value & session value - global variable
// 	session = sessionF
// 	fmt.Println("ok cookie valid, can do operation", session)

// 	//_, err = DB.Exec("UPDATE session SET resession=? WHERE user_id=?", 0, session.UserID)
// 	f(w, r)
// } else {
// 	//cookie another,  2 case -> 1 podmena, 2 antoher session active
// 	var temp string
// 	var reSession int
// 	err = DB.QueryRow("SELECT uuid FROM session WHERE user_id = ?", session.UserID).Scan(&temp)
// 	err = DB.QueryRow("SELECT resession FROM session WHERE user_id = ?", session.UserID).Scan(&reSession)
// 	if err != nil {
// 		log.Println(err, "222")
// 	}

// 	if reSession != 1 {
// 		fmt.Println("nothin")
// 	} else {
// 		fmt.Println(temp, cookie, reSession)
// 		if temp != cookie {
// 			_, err = DB.Exec("UPDATE session SET resession=? WHERE user_id=?", 1, session.UserID)
// 			utils.IsCookieExpiration(w, r, session)
// 			fmt.Println("another cookie browser, uuid != session.Uuid Db 3 ")
// 		}

// 		//				fmt.Println("signout session != 1, delete sessiob")
// 		//return
// 	}
// if err != nil {
// 	//utils.IsCookieExpiration(w, r, session)
// 	log.Println(err, "no session user in db 1")
// 	return
// }
//fmt.Println(session, "cookie exit 2", temp, cookie)

//set resession 1, when another user signin
// if resession {
// 	_, err = DB.Exec("UPDATE session SET resession=? WHERE id=?", 1, lid)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
//chaining
// func RequireAuthentication(next http.Handler) http.Handler {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			b, _ := utils.IsCookie(w, r)
// 			fmt.Print(b, "ccokie")
// 			if !b {
// 				http.Redirect(w, r, "/signin", 302)
// 				return
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// }
