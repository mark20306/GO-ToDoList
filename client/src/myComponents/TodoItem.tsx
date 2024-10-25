import { Badge, Box, Flex, Spinner, Text } from "@chakra-ui/react";
import { FaCheckCircle } from "react-icons/fa";
import { MdDelete } from "react-icons/md";
import { Todo } from "./TodoList";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { BASE_URL } from "@/App";

const TodoItem = ({ todo }: { todo: Todo }) => {
    const queryClient = useQueryClient();
    const {mutate:updateTodo, isPending:isUpdating} = useMutation({
        mutationKey:["updateTodo"],
        mutationFn: async() => {
            if(todo.completed) return alert("Todo is already completed")
                try {
                    const res = await fetch(BASE_URL + `/todos/${todo._id}`, {
                        method:"PATCH",
                    })
                    const data = await res.json()
                    if(!res.ok){
                        throw new Error(data.error || "Something went wrong")
                    }
                    return data
                } catch (error) {
                    console.log(error)
                }
        },
        onSuccess: () => {
            queryClient.invalidateQueries({queryKey:["todos"]})
        }
    });

    const { mutate: deleteTodo, isPending: isDeleting } = useMutation({
		mutationKey: ["deleteTodo"],
		mutationFn: async () => {
			try {
				const res = await fetch(BASE_URL + `/todos/${todo._id}`, {
					method: "DELETE",
				});
				const data = await res.json();
				if (!res.ok) {
					throw new Error(data.error || "Something went wrong");
				}
				return data;
			} catch (error) {
				console.log(error);
			}
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["todos"] });
		},
	});

	return (
		<Flex gap={2} alignItems={"center"} >
            <Box
                flex={1}
                borderWidth="1px"             // 使用 borderWidth 來設置外框寬度
                borderColor="gray.600"         // 外框顏色
                borderRadius="lg"              // 設置圓角
                p={2}
            >
                <Flex
                    alignItems={"center"}
                    justifyContent={"space-between"}
                >
                    <Text
                        color={todo.completed ? "green.200" : "yellow.100"}
                        textDecoration={todo.completed ? "line-through" : "none"}
                    >
                        {todo.body}
                    </Text>
                    {todo.completed && (
                        <Badge ml='1' colorPalette='green'>
                            Done
                        </Badge>
                    )}
                    {!todo.completed && (
                        <Badge ml='1' colorPalette='yellow'>
                            In Progress
                        </Badge>
                    )}
                </Flex>
            </Box>
			<Flex gap={2} alignItems={"center"}>
				<Box color={"green.500"} cursor={"pointer"} onClick={() => updateTodo()}>
                    {!isUpdating && <FaCheckCircle size={20} />}
					{isUpdating && <Spinner size={"sm"} />}
				</Box>
				<Box color={"red.500"} cursor={"pointer"} onClick={() => deleteTodo()}>
                    {!isDeleting && <MdDelete size={25} />}
                    {isDeleting && <Spinner size={"sm"} />}
				</Box>
			</Flex>
		</Flex>
	);
};
export default TodoItem;