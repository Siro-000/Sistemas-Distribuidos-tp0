import sys


def def_clientes(cantidad_clientes, f):
    for i in range(1, int(cantidad_clientes) + 1):
        f.write(f"  client{i}:\n")
        f.write(f"    container_name: client{i}\n")
        f.write("    image: client:latest\n")
        f.write("    entrypoint: /client\n")
        f.write("    networks:\n")
        f.write("      - testing_net\n")
        f.write("    depends_on:\n")
        f.write("      - server\n")
        f.write("    volumes:\n")  
        f.write(f"      - ./client/config.yaml:/config.yaml\n")
        f.write(f"      - ./.data/agency-{i}.csv:/data/agency.csv\n")
        f.write("\n")

def def_network(f):
    f.write("networks:\n")
    f.write("  testing_net:\n")
    f.write("    name: testing_net\n")
    f.write("    ipam:\n")
    f.write("      driver: default\n")
    f.write("      config:\n")
    f.write("        - subnet: 172.25.125.0/24\n")

def def_server(f, num_clients):
    f.write("  server:\n")
    f.write("    container_name: server\n")
    f.write("    image: server:latest\n")
    f.write("    entrypoint: python3 /main.py\n")
    f.write("    networks:\n")
    f.write("      - testing_net\n")
    f.write("    environment:\n")
    f.write(f"      - NUM_AGENCY={num_clients}\n") 
    f.write("    volumes:\n")  
    f.write("      - ./server/config.ini:/config.ini\n")
    f.write("\n")


def def_service(cantidad_clientes, f):
    f.write("services:\n")
        
    def_server(f, cantidad_clientes)

    def_clientes(cantidad_clientes, f)
    
    
def generar_compose(archivo_salida, cantidad_clientes):
    with open(archivo_salida, "w") as f:
        f.write("name: tp0\n")
        
        def_service(cantidad_clientes, f)

        def_network(f)

        
if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Uso: python3 mi-generador.py <archivo_salida> <cantidad_clientes>")
        sys.exit(1)

    archivo = sys.argv[1]
    cantidad = sys.argv[2]

    generar_compose(archivo, cantidad)