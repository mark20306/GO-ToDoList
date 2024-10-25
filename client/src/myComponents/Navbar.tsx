import { ColorModeButton } from "@/components/ui/color-mode";
import { Box, Flex, Text, Container } from "@chakra-ui/react";


export default function Navbar() {

	return (
		<Container maxW={"900px"}>
			
				<Flex h={16} alignItems={"center"} justifyContent={"space-between"}>
                    <Flex
						justifyContent={"center"}
						alignItems={"center"}
						gap={3}
						display={{ base: "none", sm: "flex" }}
					>
                        <img src='/go.png' alt='logo' width={40} height={40} />
						<Text fontSize={"lg"} fontWeight={500}>
							Daily Tasks
						</Text>
                        
					</Flex>

					
					<Flex alignItems={"center"} gap={3}>
						
                        <ColorModeButton />
                        
					</Flex>
				</Flex>
			
		</Container>
	);
}