import { Box, Container, Stack } from "@chakra-ui/react"
import TodoForm from "./myComponents/TodoForm"
import TodoList from "./myComponents/TodoList"
import Navbar from "./myComponents/Navbar";

export const BASE_URL = "http://localhost:5000/api";

function App() {

  return (
    <Box bg={{ base: "#8E8E8E", _dark: "#2E3447" }}>
      <Stack h='100vh' maxW="750px" marginLeft="25%">
            <Navbar />
            <Container>
              <TodoForm />
              <TodoList />
            </Container>
      </Stack>
    </Box>
	);
}

export default App
