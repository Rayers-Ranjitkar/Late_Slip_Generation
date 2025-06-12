import { Navigate } from "react-router-dom";

const ProtectedRoutes = ({children}:{children: React.ReactNode}) => { /* ReactNode bhannaley yauta component bhayo nita as children ma tw k yauta component aayeracha so, it's type is ReactNode i.e  component is one type of node in react alright */
  const isLoggedIn = false; //yo true huda matrai aba secretProtectedPage ma url ma hancha jancha i.e. children render huncha nabai return bhayera login page mai gayedincha false huda i.e. user logged in nahuda yo secret page ko access kosailai ni xaina hehe

  if(!isLoggedIn){
     //return navigate("/AdminLogin"); //if not logged in then, navigate to /Login , else render the below thing i.e. thing inside return automatically get's rendered => if not logged in then sidhai return huncha and then muni ko kura render huna paudaina +
     //useNavigate doesn't Does not return a React element but navigate to does so please import navigate so that you can do navigate to and not useNavigate //So yes Navigate to  — it returns a React element like React.createElement(Navigate, { to: "/AdminLogin", replace: true }) which tells React Router: React.createElement(Navigate, { to: "/AdminLogin", replace: true })

     return <Navigate to="/AdminLogin" />
  }

  return (
    <div>
        {children}
    </div>
  )
}

export default ProtectedRoutes