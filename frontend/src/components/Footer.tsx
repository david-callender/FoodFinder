import { FC } from "react";
import { GithubLink } from "./GithubLink";


export const Footer: FC = () => {
    return <>
        <div className="grid grid-cols-3 justify-center p-2 border-t border-gray-200 bg-black/30 py-4 text-center text-sm text-gray-200 content-center">
            
            <><p></p></>

            <footer className="relative z-10">
            
                <p>© {new Date().getFullYear()} Gopher Grub · University of Minnesota</p>
        
            </footer>
        
          
            <div className="self-center justify-self-end mr-5">
                <GithubLink />
            </div>
        
        </div>
    </>
    
}