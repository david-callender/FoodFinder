import { Footer } from "@/components/Footer";
import { LinkBar } from "@/components/LinkBar";
import { MainBackground } from "@/components/MainBackground";
import { MainPageWelcome } from "@/components/MainPageWelcome";


import type { FC } from "react";

const Home: FC = () => {
  return (
    <div className="relative flex min-h-screen flex-col overflow-hidden text-gray-900">
      <LinkBar />

      {/* Header */}
      <div className="absolute inset-0 -z-10">
        <MainBackground />
        {/* Optional dark overlay to make text pop */}
        <div className="absolute inset-0 bg-black/40 backdrop-blur-sm"></div>
      </div>

      <MainPageWelcome />


      
      <Footer />
    </div>
  );
};

export default Home;
