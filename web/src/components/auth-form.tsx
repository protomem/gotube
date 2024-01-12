import LoginForm from "./login-form";
import RegisterForm from "./register-form";
import {
  Button,
  Tab,
  TabIndicator,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
} from "@chakra-ui/react";

const AuthForm = () => {
  return (
    <Tabs isFitted variant="unstyled" w="16rem">
      <TabList borderBottomWidth="1px">
        <Tab as={Button} variant="ghost" borderRadius="0">
          Login
        </Tab>
        <Tab as={Button} variant="ghost" borderRadius="0">
          Register
        </Tab>
      </TabList>
      <TabIndicator mt="-1.5px" height="2px" bg="blue.500" borderRadius="1px" />
      <TabPanels>
        <TabPanel>
          <LoginForm />
        </TabPanel>
        <TabPanel>
          <RegisterForm />
        </TabPanel>
      </TabPanels>
    </Tabs>
  );
};

export default AuthForm;
