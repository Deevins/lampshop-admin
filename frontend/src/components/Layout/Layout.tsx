import React from "react";
import {Outlet} from "react-router-dom";
import Navbar from "../Navbar/Navbar";
import styles from "./Layout.module.scss";

const Layout: React.FC = () => {
    return (
        <div className={styles.container}>
            <Navbar/>
            <div className={styles.content}>
                <Outlet/>
            </div>
        </div>
    );
};

export default Layout;
