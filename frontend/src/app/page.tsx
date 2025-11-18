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

      <footer className="relative z-10 border-t border-white/20 bg-black/30 py-4 text-center text-sm text-gray-200">
        © {new Date().getFullYear()} Gopher Grub · University of Minnesota
      </footer>
    </div>
  );
};

export default Home;
