import { ParallaxBanner } from "react-scroll-parallax";
import { useState } from "react";
import { Banner } from "@/components/hero";
import { Info1 } from "@/components/info1";
import { Info2 } from "@/components/info2";
import { Sponsor } from "@/components/sponsors";
import { Info3 } from "@/components/info3";
import Head from "next/head";
export default function Home() {
  const [userid, setUserid] = useState("");
  const handleChane = (e) => {
    setUserid(e.target.value);
  };

  return (
    <>
      <Head>
        <title>Ingenius 2023</title>
        <link rel="shortcut icon" href="/favicon.png" />
      </Head>
      <div className="  px-4  items-center  min-h-screen py-2">
        <Banner />
        <div className="flex flex-col gap-8">
          <Info1 />
          <Info2 />
          <Info3 />
          <Sponsor />
        </div>

        <script
          async
          src="https://www.googletagmanager.com/gtag/js?id=G-Y9ZP6CZW1R"
        ></script>
        <script
          dangerouslySetInnerHTML={{
            __html: `
            window.dataLayer = window.dataLayer || [] ; function gtag()
            {dataLayer.push(arguments)}
            gtag('js', new Date()); gtag('config', 'G-Y9ZP6CZW1R');
      `,
          }}
        ></script>
      </div>

      {/* <div 
    className="flex bg-gradient-to-r from-cyan-500 to-blue-500 flex-col items-center justify-center min-h-screen py-2"
    >
        <div
        className="heading text-5xl font-lovelo-line-light lg:text-9xl font-bold text-white text-center"
        >
          <h1>INGENIUS 2023</h1>
        </div>
    </div> */}
    </>
  );
}
