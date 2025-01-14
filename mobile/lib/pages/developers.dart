import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:area/pages/home_page.dart';
import 'package:flutter/cupertino.dart';
import 'package:area/tools/screen_scale.dart';

class DevelopersPage extends StatefulWidget {
  const DevelopersPage({super.key});

  @override
  State<DevelopersPage> createState() => DevelopersPageState();
}

class DevelopersPageState extends State<DevelopersPage> {
  int currentPageIndex = 3;

  @override
  Widget build(BuildContext context) {
    final List<String> dest = [
      "/applets",
      "/create",
      "/services",
      "/plus"
    ];
    return SafeArea(
        child: Scaffold(
            bottomNavigationBar: NavigationBar(
              backgroundColor: Colors.black,
              indicatorColor: Colors.grey,
              selectedIndex: 3,
              onDestinationSelected: (int index) {
                setState(() {
                  currentPageIndex = index;
                  context.go(dest[index]);
                });
              },
              destinations: const [
                NavigationDestination(
                    icon: Icon(Icons.folder, color: Colors.white),
                    label: 'Applets'),
                NavigationDestination(
                    icon: Icon(Icons.add_circle_outline, color: Colors.white),
                    label: 'Create'),
                NavigationDestination(
                    icon: Icon(Icons.cloud, color: Colors.white),
                    label: 'Services'),
                NavigationDestination(
                    icon: Icon(CupertinoIcons.ellipsis, color: Colors.white),
                    label: 'Developers'),
              ],
            ),
            backgroundColor: Colors.white,
            body: SingleChildScrollView(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.center,
                children: [
                  const Padding(padding: EdgeInsets.all(0.0), child: MiniHeaderSection(),),
                  Align(
                    alignment: Alignment.topLeft,
                    child: ElevatedButton(
                        style: ElevatedButton.styleFrom(
                            backgroundColor: const Color.fromARGB(0, 0, 0, 0),
                            foregroundColor: const Color.fromARGB(0, 0, 0, 0),
                            shadowColor: const Color.fromARGB(0, 0, 0, 0),),
                        onPressed: () {
                          context.go("/plus");
                        },
                        child: Icon(Icons.arrow_back,
                            color: Colors.black,
                            size: screenScale(context, 0.05).height)),
                  ),
                  const Align(
                    alignment: Alignment.center,
                    child: Text(
                      'À propos de AREA',
                      style: TextStyle(
                        fontSize: 24,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  const SizedBox(height: 32),
              
                  const Text(
                    'Transformez vos actions en résultats.',
                    style:
                        TextStyle(fontSize: 18, fontWeight: FontWeight.w500),
                  ),
                  const SizedBox(height: 16),
              
                  const Text(
                    textAlign: TextAlign.center,
                    'Bienvenue sur AREA, la plateforme qui automatise vos tâches et connecte vos outils pour simplifier votre quotidien.\nQue ce soit pour optimiser votre temps, améliorer votre productivité ou gérer vos processus plus efficacement, nous sommes là pour vous accompagner.',
                    style: TextStyle(fontSize: 16),
                  ),
                  const SizedBox(height: 32),
              
                  const Text(
                    'Pourquoi AREA ?',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  _buildFeatureCard(
                    icon: Icons.connect_without_contact,
                    title: 'Connectivité illimitée',
                    description:
                        'Intégrez vos outils préférés en quelques clics.',
                  ),
                  _buildFeatureCard(
                    icon: Icons.lightbulb_outline,
                    title: 'Automatisation intelligente',
                    description:
                        'Créez des scénarios d\'actions et réactions.',
                  ),
                  _buildFeatureCard(
                    icon: Icons.settings,
                    title: 'Personnalisation complète',
                    description:
                        'Adaptez chaque fonctionnalité à vos besoins.',
                  ),
                  const SizedBox(height: 32),
              
                  const Text(
                    'Notre équipe',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  _buildTeamMember(
                    name: 'Paul',
                    role: 'Développeur Back-End',
                    description:
                        'Expert en infrastructures, il garantit la performance et la fiabilité de la plateforme.',
                  ),
                  _buildTeamMember(
                    name: 'Grégoire',
                    role: 'Product Owner & Développeur Full Stack',
                    description:
                        'Le guide stratégique qui allie vision et développement pour concrétiser le projet.',
                  ),
                  _buildTeamMember(
                    name: 'Élise',
                    role: 'Développeuse Front Mobile',
                    description:
                        'Elle conçoit des interfaces mobiles élégantes et intuitives pour une expérience utilisateur fluide.',
                  ),
                  _buildTeamMember(
                    name: 'John',
                    role: 'Responsable DevOps',
                    description:
                        'Il s’assure de la stabilité et de l’évolutivité des systèmes pour un fonctionnement sans faille.',
                  ),
                  _buildTeamMember(
                    name: 'Esteban',
                    role: 'Développeur Front Web',
                    description:
                        'Créateur d’interfaces web modernes, il rend chaque interaction agréable et efficace.',
                  ),
                  const SizedBox(height: 32),
              
                  const Text(
                    'Notre mission et nos valeurs',
                    style: TextStyle(
                      fontSize: 20,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 16),
                  const Text(
                    'Chez AREA, nous croyons que l’innovation peut simplifier la vie. Notre mission est de :',
                    style: TextStyle(fontSize: 16),
                  ),
                  const SizedBox(height: 8),
                  const Text(
                    '- Automatiser les tâches quotidiennes pour vous faire gagner du temps.\n'
                    '- Rendre la technologie accessible à travers une interface intuitive.\n'
                    '- Favoriser l’efficacité et la créativité, en laissant la machine gérer le répétitif.',
                    style: TextStyle(fontSize: 16),
                  ),
                ],
              ),
            )));
  }

  Widget _buildFeatureCard(
      {required IconData icon,
      required String title,
      required String description}) {
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      child: ListTile(
        leading: Icon(icon, size: 40, color: Colors.blueAccent),
        title: Text(title,
            style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w600)),
        subtitle: Text(description, style: const TextStyle(fontSize: 16)),
      ),
    );
  }

  Widget _buildTeamMember(
      {required String name,
      required String role,
      required String description}) {
    return Card(
      margin: const EdgeInsets.only(bottom: 16),
      child: ListTile(
        title: Text('$name – $role',
            style: const TextStyle(fontSize: 18, fontWeight: FontWeight.w600)),
        subtitle: Text(description, style: const TextStyle(fontSize: 16)),
      ),
    );
  }
}
