import Head from "next/head";
import React, { useState } from "react";
import BarcodeScannerComponent from "@steima/react-qr-barcode-scanner";
import Image from "next/image";
import { Inter } from "next/font/google";
import axios from "axios";
import 'bootstrap/dist/css/bootstrap.css';
import styles from "@/styles/Home.module.css";

const inter = Inter({ subsets: ["latin"] });
const API_URL = "http://localhost:5001";

export default function Home() {
   const [data, setData] = React.useState("Not Found");
   const [details, setDetails] = React.useState(true);
   const [data1, setData1] = React.useState();
   const [attempt, setAttempt] = React.useState(false)
   const [qrscan, setqrscan] = React.useState()
   const [qr, setqr] = React.useState(false)
   //...
   const [query, setQuery] = useState({
      SRN: "",
      name: "",
      phone: "",
      email: "",
   });
   const handleChange = () => (e) => {
      const name = e.target.name;
      const value = e.target.value;
      setQuery((prevState) => ({
         ...prevState,
         [name]: value,
      }));
   };

   const handleSubmit = (e) => {
      e.preventDefault();
      const formData = new FormData();
      Object.entries(query).forEach(([key, value]) => {
         formData.append(key, value);
      });
      const jsonObject = {};

      for (const [key, value] of formData) {
         jsonObject[key] = value;
      }

      // console.log(query);

      axios
         .post(`${API_URL}/checkin`, jsonObject, {
            headers: { Accept: "application/json" },
         })
         .then(function(response) {
            setAttempt(false)
            setAttempt(true)
            setQuery({
               SRN: "",
               name: "",
               phone: "",
               email: "",
            });
            if (response.data.status == true) {
               setData("Found user. " + response.data.message)
               setqrscan(response.data.QR);
               setqr(true)
            } else {
               setData("User records not found :( " + response.data.message)
               setqrscan("");
            }
         })
         .catch(function(error) {
            setAttempt(false)
            setAttempt(true)
            setData("User records not found :( ")
            setqr(false)
            console.log(error);
         });
   };

   return (
      <>
         <Head>
            <title>Participant Check-in</title>
            <meta name="description" content="Generated by create next app" />
            <meta name="viewport" content="width=device-width, initial-scale=1" />
            <link rel="icon" href="/favicon.ico" />
         </Head>
         <main className={styles.main}>
            <h1> Participant Check-in </h1>
            <div className="form-group">
               <form
                  acceptCharset="UTF-8"
                  method="POST"
                  id="ajaxForm"
                  onSubmit={handleSubmit}
               >

                  <div className="form-group">
                     <label>SRN : </label>
                     <input
                        type="text"
                        name="SRN"
                        className="form-control"
                        value={query.SRN}
                        onChange={handleChange()}
                        required
                     ></input>

                  </div>
                  <br />
                  <div className="form-group">
                     <label>Name : </label>
                     <input
                        type="text"
                        name="name"
                        className="form-control"
                        value={query.name}
                        onChange={handleChange()}
                        required
                     ></input>

                  </div>                  <br />
                  <div className="form-group">
                     <label>Email : </label>
                     <input
                        type="email"
                        name="email"
                        className="form-control"
                        value={query.email}
                        onChange={handleChange()}
                        required
                     ></input>

                  </div>                     <br />
                  <div className="form-group">
                     <label>Phone : </label>
                     <input
                        type="text"
                        name="phone"
                        className="form-control"
                        value={query.phone}
                        onChange={handleChange()}
                        required
                     ></input>
                  </div>
                  <button type="submit" className="btn btn-primary mt-4 position-absolute start-50 translate-middle">
                     Submit
                  </button>
               </form>
            </div>
            {attempt ? (
               <div class="d-flex align-items-center justify-content-center flex-column" >
                  <p>{data}</p>
                  {qr ? <div class="d-flex align-items-center justify-content-center flex-column">
                     <p>Click on the qr code to download it as an image</p> <br />
                     <a href={"data:image/png;base64," + qrscan} download="myimage">
                        <Image
                           id="imgElem"
                           src={"data:image/png;base64," + qrscan}
                           alt="Picture of the author"
                           width={200}
                           height={200}
                        ></Image>
                     </a> </div> : null}
               </div>
            ) : <p>Please enter your details!</p>}
            <p> Please enter same details as entered during registration, if you feel that the check in system is not working as intended please reach out to any of our members from organising comiteee </p>
         </main>
      </>
   );
}
