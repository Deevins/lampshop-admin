import React from "react";
import {NavLink, useNavigate} from "react-router-dom";
import styles from "./Navbar.module.scss";
import {logout} from "../../api/AuthApi.ts";

const Navbar: React.FC = () => {
    const navigate = useNavigate();

    const handleLogout = () => {
        logout();
        navigate("/login");
    };


    return (
        <div className={styles.navbar}>
            <div className={styles["nav-links"]}>
                <NavLink
                    to="/products"
                    className={({isActive}) =>
                        isActive ? `${styles["nav-link"]} ${styles["active"]}` : styles["nav-link"]
                    }
                >
                    Товары
                </NavLink>
                <NavLink
                    to="/orders"
                    className={({isActive}) =>
                        isActive ? `${styles["nav-link"]} ${styles["active"]}` : styles["nav-link"]
                    }
                >
                    Заказы
                </NavLink>
            </div>
            <button className={styles["logout-button"]} onClick={handleLogout}>
                Выход
            </button>
        </div>
    );
};

export default Navbar;
